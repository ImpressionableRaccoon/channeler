package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	envTelegramAppID             = "TELEGRAM_APP_ID"
	envTelegramAppHash           = "TELEGRAM_APP_HASH"
	envSessionStoragePath        = "SESSION_STORAGE_PATH"
	envTelegramChannelID         = "TELEGRAM_CHANNEL_ID"
	envTelegramChannelAccessHash = "TELEGRAM_CHANNEL_ACCESS_HASH"
	envYDBConnectionString       = "YDB_CONNECTION_STRING"
	envTablePathPrefix           = "TABLE_PATH_PREFIX"
)

type config struct {
	TelegramAppID             int
	TelegramAppHash           string
	SessionStoragePath        string
	TelegramChannelID         int64
	TelegramChannelAccessHash int64
	YDBConnectionString       string
	TablePathPrefix           string
}

func Load() (config, error) {
	c := config{}

	var err error
	var exists bool

	c.TelegramAppID, err = strconv.Atoi(os.Getenv(envTelegramAppID))
	if err != nil {
		return config{}, fmt.Errorf("%w: %s: %s", ErrKeyParse, envTelegramAppID, err.Error())
	}

	c.TelegramAppHash, exists = os.LookupEnv(envTelegramAppHash)
	if !exists {
		return config{}, fmt.Errorf("%w: %s", ErrKeyNotExists, envTelegramAppHash)
	}

	c.SessionStoragePath = os.Getenv(envSessionStoragePath)

	c.TelegramChannelID, err = strconv.ParseInt(os.Getenv(envTelegramChannelID), 10, 64)
	if err != nil {
		return config{}, fmt.Errorf("%w: %s: %s", ErrKeyParse, envTelegramChannelID, err.Error())
	}

	c.TelegramChannelAccessHash, err = strconv.ParseInt(os.Getenv(envTelegramChannelAccessHash), 10, 64)
	if err != nil {
		return config{}, fmt.Errorf("%w: %s: %s", ErrKeyParse, envTelegramChannelAccessHash, err.Error())
	}

	c.YDBConnectionString = os.Getenv(envYDBConnectionString)
	c.TablePathPrefix = os.Getenv(envTablePathPrefix)

	return c, nil
}
