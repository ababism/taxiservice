package mshttp

import (
	"time"
)

const (
	defaultHost         = "localhost"
	defaultPort         = "8080"
	defaultReadTimeout  = time.Second
	defaultWriteTimeout = time.Second
)

type Config struct {
	Host         string        `mapstructure:"host" yaml:"host" env:"HOST"`
	Port         string        `mapstructure:"port" yaml:"port" env:"PORT"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" yaml:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" yaml:"write_timeout" env:"WRITE_TIMEOUT"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Host:         defaultHost,
		Port:         defaultPort,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}
}
