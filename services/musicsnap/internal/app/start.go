package app

import (
	"context"
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	requestid "github.com/sumit-tembe/gin-requestid"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"music-snap/pkg/msshutdown"
	"time"
	//"music-snap/pkgmetrics"
	httpServer "music-snap/pkg/mshttp"
	myHttp "music-snap/services/musicsnap/internal/handler/http"
	oapigen "music-snap/services/musicsnap/internal/handler/http/musicsnap/oapi"
)

// Start - Единая точка запуска приложения
func (a *App) Start(ctx context.Context) {

	_, err := a.cfg.CacheRefresher.GetIterationInterval()
	if err != nil {
		a.logger.Fatal("can't parse time from daemon config string:", zap.Error(err))
	}

	//a.daemon.Start(daemonInterval)

	go a.startHTTPServer(ctx)

	if err := msshutdown.Wait(a.cfg.GracefulShutdown); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to gracefully shutdown %s app: %s", a.cfg.App.Name, err.Error()))
	} else {
		a.logger.Info("App gracefully stopped")
	}
}

func (a *App) startHTTPServer(ctx context.Context) {
	// Создаем общий роутинг http сервера
	router := httpServer.NewRouter()
	//
	//// Добавляем системные роуты
	//router.WithHandleGET("/metrics", metrics.HandleFunc())

	tracerMw := oapigen.MiddlewareFunc(otelgin.Middleware(a.cfg.App.Name, otelgin.WithTracerProvider(a.tracerProvider)))
	GinZapMw := oapigen.MiddlewareFunc(ginzap.Ginzap(a.logger, time.RFC3339, true))
	requestIdMw := oapigen.MiddlewareFunc(requestid.RequestID(nil))
	middlewares := []oapigen.MiddlewareFunc{
		tracerMw,
		GinZapMw,
		requestIdMw,
	}

	// Добавляем роуты api
	myHttp.InitHandler(router.Router(), a.logger, middlewares, a.service)

	// Создаем сервер
	srv := httpServer.NewServer(a.cfg.Http, router)
	//srv.RegisterRoutes(&router)

	// Стартуем
	a.logger.Info(fmt.Sprintf("Starting %s HTTP server at %s:%s", a.cfg.App.Name, a.cfg.Http.Host, a.cfg.Http.Port))
	if err := srv.Start(); err != nil {
		a.logger.Error(fmt.Sprintf("Fail with %s HTTP server: %s", a.cfg.App.Name, err.Error()))
		msshutdown.Now()
	}
}
