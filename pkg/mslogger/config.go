package mslogger

type Config struct {
	Level        string   `mapstructure:"level"`
	Env          string   `mapstructure:"env"`
	Outputs      []string `mapstructure:"outputs"`
	ErrorOutputs []string `mapstructure:"error_outputs"`
	Encoding     string   `mapstructure:"encoding"`
	SentryLevel  string   `mapstructure:"sentry_level"`
	SentryDSN    string   `mapstructure:"sentry_dsn"`
}
