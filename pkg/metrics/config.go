package metrics

type Config struct {
	MetricsEnable  bool   `mapstructure:"enable"`
	MetricsAddress string `mapstructure:"address"`
}
