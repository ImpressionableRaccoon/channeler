package config

import "errors"

var (
	ErrKeyNotExists = errors.New("key not exists")
	ErrKeyParse     = errors.New("key parse error")
)
