package postgre

type Config struct {
	DbDSN string `mapstructure:"db_dsn" env:"db_dsn"`
}
