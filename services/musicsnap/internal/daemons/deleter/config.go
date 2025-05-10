package deleter

import "time"

type Config struct {
	IterationInterval string `mapstructure:"iteration_interval"`
}

func (c Config) GetIterationInterval() (time.Duration, error) {
	return time.ParseDuration(c.IterationInterval)
}
