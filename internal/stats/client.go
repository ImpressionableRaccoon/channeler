package stats

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

type client struct {
	t        *telegram.Client
	stopFunc bg.StopFunc
}

var _ StatManager = client{}

func New(logger *zap.Logger, sessionStoragePath string, appID int, appHash string) (client, error) {
	c := client{}

	opts := telegram.Options{
		Logger: logger,
		Device: telegram.DeviceConfig{
			DeviceModel: "Channeler",
		},
	}

	if len(sessionStoragePath) > 0 {
		opts.SessionStorage = &telegram.FileSessionStorage{
			Path: sessionStoragePath,
		}
		logger.Info("session storage path specified", zap.String("path", sessionStoragePath))
	} else {
		logger.Warn("session storage path not specified")
	}

	c.t = telegram.NewClient(appID, appHash, opts)

	var err error
	c.stopFunc, err = bg.Connect(c.t)
	if err != nil {
		return client{}, fmt.Errorf("telegram connect failed: %w", err)
	}

	s, err := c.t.Auth().Status(context.Background())
	if err != nil {
		return client{}, fmt.Errorf("error getting auth status: %w", err)
	}
	if !s.Authorized {
		err = c.authorizeUser()
		if err != nil {
			return client{}, fmt.Errorf("error authorizing user: %w", err)
		}
		s, err = c.t.Auth().Status(context.Background())
		if err != nil {
			return client{}, fmt.Errorf("error getting auth status: %w", err)
		}
	}

	logger.Info("authorized", zap.Any("authorization info", s))

	return c, nil
}

func (c client) Stop() error {
	return c.stopFunc()
}

func (c client) authorizeUser() error {
	readString := func(prompt string) (string, error) {
		fmt.Printf("%s: ", prompt)
		res, err := bufio.NewReader(os.Stdin).ReadString('\n')
		return strings.TrimSpace(res), err
	}

	phoneNumber, err := readString("Enter phone number")
	if err != nil {
		return err
	}
	password, err := readString("Enter password")
	if err != nil {
		return err
	}

	return auth.NewFlow(
		auth.Constant(phoneNumber, password,
			auth.CodeAuthenticatorFunc(
				func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
					return readString("Enter code")
				},
			),
		),
		auth.SendCodeOptions{},
	).Run(context.Background(), c.t.Auth())
}
