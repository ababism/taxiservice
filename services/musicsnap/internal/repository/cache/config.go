package cache

import "time"

type Config struct {
	Expiration  time.Duration `mapstructure:"expiration"`
	InitialSize int           `mapstructure:"initial_size"`
}
