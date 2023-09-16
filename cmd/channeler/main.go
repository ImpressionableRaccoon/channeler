package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ImpressionableRaccoon/channeler/internal/config"
	"github.com/ImpressionableRaccoon/channeler/internal/stats"
)

var (
	logger *zap.Logger
	err    error
)

func init() {
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

	t, err := stats.New(logger, c.SessionStoragePath, c.TelegramAppID, c.TelegramAppHash)
	if err != nil {
		logger.Panic("init stats client failed", zap.Error(err))
	}
	defer func() {
		err := t.Stop()
		if err != nil {
			logger.Error("error stopping stats client", zap.Error(err))
		}
	}()

	<-ctx.Done()
}
