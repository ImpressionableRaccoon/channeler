package ydb

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"time"

	environ "github.com/ydb-platform/ydb-go-sdk-auth-environ"
	ydbZap "github.com/ydb-platform/ydb-go-sdk-zap"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/balancers"
	"github.com/ydb-platform/ydb-go-sdk/v3/trace"
	"go.uber.org/zap"

	"github.com/ImpressionableRaccoon/channeler/internal/datarealm"
)

type ydbStorage struct {
	logger *zap.Logger

	cc *ydb.Driver
	c  ydb.SQLConnector
	db *sql.DB
}

var _ datarealm.Storager = ydbStorage{}

func New(ctx context.Context, logger *zap.Logger, dsn string, pathPrefix string) (*ydbStorage, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	s := &ydbStorage{
		logger: logger.Named("ydbStorage"),
	}
	var err error

	s.cc, err = ydb.Open(ctx,
		dsn,
		ydb.WithBalancer(balancers.SingleConn()),
		environ.WithEnvironCredentials(ctx),
		ydbZap.WithTraces(
			logger,
			trace.DetailsAll,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("ydb.Open error: %w", err)
	}

	pathPrefix = path.Join(s.cc.Name(), pathPrefix)

	s.c, err = ydb.Connector(s.cc,
		ydb.WithAutoDeclare(),
		ydb.WithTablePathPrefix(pathPrefix),
	)
	if err != nil {
		_ = s.cc.Close(ctx)
		return nil, fmt.Errorf("ydb.Connector error: %w", err)
	}

	s.db = sql.OpenDB(s.c)

	err = s.doMigrations(ctx)
	if err != nil {
		_ = s.db.Close()
		_ = s.c.Close()
		_ = s.cc.Close(ctx)
		return nil, fmt.Errorf("migrations error: %w", err)
	}

	return s, nil
}

func (s ydbStorage) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("error close database/sql driver: %w", err)
	}

	if err := s.c.Close(); err != nil {
		return fmt.Errorf("error close connector: %w", err)
	}

	if err := s.cc.Close(ctx); err != nil {
		return fmt.Errorf("error close ydb driver: %w", err)
	}

	return nil
}
