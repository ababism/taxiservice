package app

import "fmt"

type Config struct {
	Name        string      `mapstructure:"name"`
	Environment Environment `mapstructure:"env"`
	Version     string      `mapstructure:"version"`
	Dc          string      `mapstructure:"dc"`
}

type RawConfig struct {
	Service     string `yaml:"service" env:"SERVICE"`
	Environment string `yaml:"environment" env:"ENVIRONMENT"`
	Dc          string `yaml:"dc" env:"DC"`
}

func (r *RawConfig) ToConfig() (*Config, error) {
	parsedEnv, err := ParseEnvironment(r.Environment)
	if err != nil {
		return nil, err
	}
	return &Config{
		Name:        r.Service,
		Environment: parsedEnv,
		Dc:          r.Dc,
	}, nil
}

func (c *Config) IsProduction() bool {
	return c.Environment == ProductionEnv
}

func (c *Config) IsTesting() bool {
	return c.Environment == TestingEnv
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == DevelopmentEnv
}

func (c *Config) Release() string {
	return fmt.Sprintf("%s-%s", c.Environment, c.Version)
}
