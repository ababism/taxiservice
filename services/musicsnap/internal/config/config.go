package config

import (
	"github.com/spf13/viper"
	"log"
	"music-snap/pkg/app"
	msconfig "music-snap/pkg/config"
	"music-snap/pkg/metrics"
	"music-snap/pkg/msdb/mspostgres"
	"music-snap/pkg/mshttp"
	"music-snap/pkg/mslogger"
	"music-snap/pkg/msshutdown"
	"music-snap/pkg/mstracer"
	"music-snap/services/musicsnap/internal/daemons/cacherefresher"
	"music-snap/services/musicsnap/internal/repository/cache"
	"music-snap/services/musicsnap/internal/service/jwtservice"
	//"music-snap/services/musicsnap/internal/repository/postgre"
)

type Config struct {
	App              *app.Config            `mapstructure:"app"`
	Http             *mshttp.Config         `mapstructure:"http"`
	Logger           *mslogger.Config       `mapstructure:"logger"`
	Metrics          *metrics.Config        `mapstructure:"metrics"`
	GracefulShutdown *msshutdown.Config     `mapstructure:"graceful_shutdown"`
	Tracer           *mstracer.Config       `mapstructure:"tracer"`
	CacheRefresher   *cacherefresher.Config `mapstructure:"cache_refresher"`
	Cache            *cache.Config          `mapstructure:"cache"`
	Postgres         *mspostgres.Config     `mapstructure:"postgres"`
	JWTService       *jwtservice.Config     `mapstructure:"jwtservice"`
}

func NewConfig(filePath string, appName string) (*Config, error) {
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error while reading config file: %v", err)
	}

	// Загрузка конфигурации в структуру Config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("error while unmarshalling config file: %v", err)
	}

	// Замена значений из переменных окружения, если они заданы
	msconfig.ReplaceWithEnv(&config, appName)
	return &config, nil
}
