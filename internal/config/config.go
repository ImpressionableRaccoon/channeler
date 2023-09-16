package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	envTelegramAppID      = "TELEGRAM_APP_ID"
	envTelegramAppHash    = "TELEGRAM_APP_HASH"
	envSessionStoragePath = "SESSION_STORAGE_PATH"
)

type config struct {
	TelegramAppID      int
	TelegramAppHash    string
	SessionStoragePath string
}

func Load() (config, error) {
	c := config{}
	var err error

	var value string
	var exists bool
	value, exists = os.LookupEnv(envTelegramAppID)
	if !exists {
		return config{}, fmt.Errorf("%w: %s", ErrKeyNotExists, envTelegramAppID)
	}
	c.TelegramAppID, err = strconv.Atoi(value)
	if err != nil {
		return config{}, fmt.Errorf("%w: %s: %s", ErrKeyParse, envTelegramAppID, err.Error())
	}

	c.TelegramAppHash, exists = os.LookupEnv(envTelegramAppHash)
	if !exists {
		return config{}, fmt.Errorf("%w: %s", ErrKeyNotExists, envTelegramAppHash)
	}

	c.SessionStoragePath = os.Getenv(envSessionStoragePath)

	return c, nil
}
