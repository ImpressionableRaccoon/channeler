package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ImpressionableRaccoon/channeler/internal/config"
	"github.com/ImpressionableRaccoon/channeler/internal/stats"
	"github.com/ImpressionableRaccoon/channeler/internal/storage/ydb"
	"github.com/ImpressionableRaccoon/channeler/internal/workers"
)

var logger *zap.Logger

func init() {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	var err error
	logger, err = cfg.Build(zap.AddStacktrace(zapcore.PanicLevel))
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

	w, err := workers.New(logger, st, &sm)
	if err != nil {
		logger.Panic("init workers failed", zap.Error(err))
	}

	go w.AdminLogUpdater(ctx, time.Minute, c.TelegramChannelID, c.TelegramChannelAccessHash)

	<-ctx.Done()
}
