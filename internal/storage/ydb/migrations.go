package ydb

import (
	"context"
	"database/sql"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/retry"
)

const (
	createEventsTableQuery = `
CREATE TABLE events (
    id          Int64 NOT NULL,
    date        Timestamp,
    user_id     Int64,
    action_type Utf8,
    action      Json,
    PRIMARY KEY (
                id
    )
);
`
	createUsersTableQuery = `
CREATE TABLE users (
    id          Int64 NOT NULL,
    access_hash Int64,
	first_name  Utf8,
    last_name   Utf8,
    username    Utf8,
    phone       Utf8,
    PRIMARY KEY (
                id
    )
);
`
)

func (s ydbStorage) doMigrations(ctx context.Context) error {
	var err error
	if err = ctx.Err(); err != nil {
		return err
	}

	err = s.createTable(ctx, createEventsTableQuery)
	if err != nil {
		return err
	}

	err = s.createTable(ctx, createUsersTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func (s ydbStorage) createTable(ctx context.Context, query string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return retry.Do(ydb.WithTxControl(ctx, writeTx), s.db,
		func(ctx context.Context, cc *sql.Conn) error {
			_, err := s.db.ExecContext(ydb.WithQueryMode(ctx, ydb.SchemeQueryMode), query)
			return err
		}, retry.WithDoRetryOptions(retry.WithIdempotent(true)),
	)
}
