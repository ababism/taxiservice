package app

import (
	"fmt"
	"strings"
)

const (
	DevelopmentEnv     = "dev"
	TestingEnv         = "test"
	ProductionEnv      = "prod"
	UnknownEnvironment = "unknown"
)

type Environment string

func ParseEnvironment(env string) (Environment, error) {
	switch strings.ToLower(env) {
	case DevelopmentEnv:
		return DevelopmentEnv, nil
	case TestingEnv:
		return TestingEnv, nil
	case ProductionEnv:
		return ProductionEnv, nil
	default:
		return UnknownEnvironment, fmt.Errorf("got unsupported environment: %s", env)
	}
}
