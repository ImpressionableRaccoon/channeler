package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ImpressionableRaccoon/channeler/internal/config"
	"github.com/ImpressionableRaccoon/channeler/internal/datarealm"
	"github.com/ImpressionableRaccoon/channeler/internal/stats"
	"github.com/ImpressionableRaccoon/channeler/internal/storage/ydb"
)

const filePath = "people.csv"

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		panic(fmt.Errorf("error create logger: %w", err))
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	defer func() { _ = logger.Sync() }()

	c, err := config.Load()
	if err != nil {
		logger.Panic("load config failed", zap.Error(err))
	}

	st, err := ydb.New(ctx, logger, c.YDBConnectionString, c.TablePathPrefix)
	if err != nil {
		logger.Panic("init ydb storager failed", zap.Error(err))
	}
	defer func() {
		err := st.Stop(context.Background())
		if err != nil {
			logger.Error("error stopping ydb storage", zap.Error(err))
		}
	}()

	sm, err := stats.New(logger, c.SessionStoragePath, c.TelegramAppID, c.TelegramAppHash)
	if err != nil {
		logger.Panic("init stats client failed", zap.Error(err))
	}
	defer func() {
		err := sm.Stop()
		if err != nil {
			logger.Error("error stopping stats client", zap.Error(err))
		}
	}()

	f, err := os.Open(filePath)
	if err != nil {
		logger.Panic("error opening file",
			zap.String("file path", filePath),
			zap.Error(err),
		)
	}
	defer func() { _ = f.Close() }()

	csvReader := csv.NewReader(f)
	people, err := csvReader.ReadAll()
	if err != nil {
		logger.Panic("error reading file", zap.Error(err))
	}

	for i, person := range people {
		var userID, userAccessHash int64
		var user datarealm.User

		userID, err = strconv.ParseInt(person[0], 10, 64)
		if err != nil {
			logger.Error("parsing user ID failed",
				zap.Int("line", i),
				zap.Error(err),
			)
			continue
		}
		userAccessHash, err = strconv.ParseInt(person[1], 10, 64)
		if err != nil {
			logger.Error("parsing user access hash failed",
				zap.Int("line", i),
				zap.Error(err),
			)
			continue
		}

		// sleep to prevent FLOOD_WAIT error
		time.Sleep(5 * time.Second)

		user, err = sm.GetUser(ctx, userID, userAccessHash)
		if err != nil {
			logger.Error("get user failed",
				zap.Int64("user ID", userID),
				zap.Int64("access hash", userAccessHash),
				zap.Error(err),
			)
			continue
		}

		err = st.AddUser(ctx, user)
		if err != nil {
			logger.Error("add user failed",
				zap.Any("user", user),
				zap.Error(err),
			)
		}
	}
}
