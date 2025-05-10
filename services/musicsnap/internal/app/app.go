package app

import (
	"context"
	"fmt"
	"github.com/juju/zaputil/zapctx"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"music-snap/pkg/app"
	"music-snap/pkg/metrics"
	"music-snap/pkg/msdb/mspostgres"
	"music-snap/pkg/mslogger"
	"music-snap/pkg/msshutdown"
	"music-snap/pkg/mstracer"
	"music-snap/services/musicsnap/internal/config"
	"music-snap/services/musicsnap/internal/daemons/cacherefresher"
	"music-snap/services/musicsnap/internal/repository/cache"
	"music-snap/services/musicsnap/internal/repository/postgre"
	"music-snap/services/musicsnap/internal/service"
	"music-snap/services/musicsnap/internal/service/jwtservice"
	//"music-snap/services/musicsnap/internal/service/ports"
)

type App struct {
	cfg            *config.Config
	address        string
	logger         *zap.Logger
	tracerProvider *trace.TracerProvider
	service        service.MusicSnapService
	daemon         *cacherefresher.CacheRefresher
}

func NewApp(cfg *config.Config) (*App, error) {
	startCtx := context.Background()
	// INFRASTRUCTURE ----------------------------------------------------------------------

	// Инициализируем logger
	//logger, err := xlogger.Init(cfg.Logger, cfg.App)
	logger, err := mslogger.InitLogger(cfg.Logger, cfg.App.Name)
	if err != nil {
		return nil, err
	}
	// Чистим кэш logger при shutdown
	msshutdown.AddCallback(
		&msshutdown.Callback{
			Name: "ZapLoggerCacheWipe",
			FnCtx: func(ctx context.Context) error {
				return logger.Sync()
			},
		})
	logger.Info("Init Logger – success")

	logger.Debug("Logger Debug test")
	logger.Info("Logger Info test")
	logger.Error("Logger Error test")

	// сохраняем logger в контекст
	_ = zapctx.WithLogger(startCtx, logger)

	// Инициализируем обработку ошибок
	err = app.InitAppError(cfg.App)
	if err != nil {
		logger.Fatal("while initializing App Error handling package", zap.Error(err))
	}

	logger.Info("initializing App Error – success")

	logger.Info("App environment: ", zap.String("env", string(cfg.App.Environment)))

	// Инициализируем трассировку
	tp, logFileCloseOptional, err := mstracer.Init(cfg.Tracer, cfg.App)
	if err != nil {
		return nil, err
	}
	msshutdown.AddCallback(
		&msshutdown.Callback{
			Name: "OpenTelemetryShutdown",
			FnCtx: func(ctx context.Context) error {
				if err := tp.Shutdown(context.Background()); err != nil {
					logger.Error("Error shutting down tracer provider: %v", zap.Error(err))
					return err
				}
				return nil
			},
		})
	if logFileCloseOptional != nil {
		msshutdown.AddCallback(
			&msshutdown.Callback{
				Name: "OTelLogFileClosure",
				FnCtx: func(ctx context.Context) error {
					if err := logFileCloseOptional(); err != nil {
						logger.Error("Error closing down tracer log file: %v", zap.Error(err))
						return err
					}
					return nil
				},
			})
	}
	logger.Info("Init Tracer – success")

	// Инициализируем Prometheus
	metrics.InitOnce(cfg.Metrics, logger, metrics.AppInfo{
		Name:        cfg.App.Name,
		Environment: string(cfg.App.Environment),
		Version:     cfg.App.Version,
	})
	logger.Info("Init Metrics – success")

	// Инициализируем вспомогательные сервисы
	_, err = jwtservice.New(cfg.JWTService)
	if err != nil {
		logger.Fatal("Error init JWTService:", zap.Error(err))
		return nil, errors.Wrap(err, "Init JWTService")
	}

	// REPOSITORY ----------------------------------------------------------------------

	// Инициализация PostgreSQL
	PostgreSQL, PsqlClose, err := mspostgres.NewDB(cfg.Postgres)
	if err != nil {
		logger.Fatal("Error init Postgres DB:", zap.Error(err))
		return nil, errors.Wrap(err, "Init Postgres DB")
	}
	msshutdown.AddCallback(
		&msshutdown.Callback{
			Name: "PSQLClientDisconnect",
			FnCtx: func(ctx context.Context) error {
				return PsqlClose()
			},
		})
	logger.Info("PosgtreSQL connect – success")
	logger.Info(fmt.Sprintf("PosgtreSQL driver: %s", PostgreSQL.DriverName()))

	// Инициализируем кэш
	profileCache := cache.New(cfg.Cache)

	// All repositories
	jwtService, err := jwtservice.New(cfg.JWTService)
	if err != nil {
		logger.Fatal("Error init JWTService:", zap.Error(err))
		return nil, errors.Wrap(err, "Init JWTService")
	}

	repos := postgre.NewRepository(PostgreSQL)

	// User
	//userRepository := postgre.NewUserRepository(PostgreSQL)

	// SERVICE LAYER ----------------------------------------------------------------------

	// Service layer

	//bannerService := service.NewBannerService(bannerRepository, profileCache)
	musicSnapService := service.New(repos, jwtService, profileCache)

	//userSvc := service.NewUserSvc(userRepository, jwtService, profileCache)
	//authSvc := service.NewAuthSvc(jwtService, userRepository)
	//subSvc := service.NewSubscriptionSvc(userRepository, profileCache)
	//
	//musicSnapService := service.MusicSnapService{
	//	Notification: nil,
	//	Auth:         authSvc,
	//	User:         userSvc,
	//	Subscription: subSvc,
	//	Review:       nil,
	//	Reaction:     nil,
	//	Photo:        nil,
	//	Stats:        nil,
	//	Event:        nil,
	//	Note:         nil,
	//	Playlist:     nil,
	//	//Banner:       bannerService,
	//}

	logger.Info(fmt.Sprintf("Init %s – success", cfg.App.Name))

	//CacheRefresher for cache refreshing
	daemon := cacherefresher.New(logger, profileCache)
	msshutdown.AddCallback(
		&msshutdown.Callback{
			Name:  "cache refresher daemon stop",
			FnCtx: daemon.StopFunc(),
		})

	daemonInterval, err := cfg.CacheRefresher.GetIterationInterval()
	if err != nil {
		logger.Fatal("can't parse time from daemon config string:", zap.Error(err))
	}
	logger.Info("CacheRefresher interval – ", zap.Duration("interval", daemonInterval))

	logger.Info("Init CacheRefresher – success")

	//service.NewMusicSnapService()

	// TRANSPORT LAYER ----------------------------------------------------------------------

	// инициализируем адрес сервера
	address := fmt.Sprintf(":%d", cfg.Http.Port)

	return &App{
		cfg:            cfg,
		logger:         logger,
		service:        musicSnapService,
		address:        address,
		tracerProvider: tp,
		daemon:         daemon,
	}, nil
}
