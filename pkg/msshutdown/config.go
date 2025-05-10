package msshutdown

import (
	"time"
)

const (
	defaultDelay           = 5 * time.Second
	defaultWaitTimeout     = 10 * time.Second
	defaultCallbackTimeout = 2 * time.Second
)

type Config struct {
	Delay           time.Duration `mapstructure:"delay"`
	WaitTimeout     time.Duration `mapstructure:"wait_timeout"`
	CallbackTimeout time.Duration `mapstructure:"callback_timeout"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Delay:           defaultDelay,
		WaitTimeout:     defaultWaitTimeout,
		CallbackTimeout: defaultCallbackTimeout,
	}
}
