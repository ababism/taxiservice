package cacherefresher

import (
	"context"
	"github.com/google/uuid"
	"github.com/juju/zaputil/zapctx"
	global "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/domain/keys"
	"music-snap/services/musicsnap/internal/service/ports"
	"runtime"
	"time"
)

type CacheRefresher struct {
	started   bool
	stop      chan bool
	userCache ports.ProfileCache
	logger    *zap.Logger
}

func New(logger *zap.Logger, cache ports.ProfileCache) *CacheRefresher {
	return &CacheRefresher{
		logger:    logger,
		userCache: cache,
		stop:      make(chan bool),
		started:   false}
}

func (s *CacheRefresher) stopCallback(ctx context.Context) error {
	if s.started != true {
		return nil
	}
	s.started = false
	s.stop <- true
	return nil
}

func (s *CacheRefresher) StopFunc() func(context.Context) error {
	return s.stopCallback
}

func (s *CacheRefresher) Start(scrapeInterval time.Duration) {
	s.started = true
	go func() {
		stop := s.stop
		go func() {
			for {
				s.refresh(scrapeInterval)
				runtime.Gosched()
			}
		}()
		<-stop
	}()
}

func generateRequestID() string {
	id := uuid.New()
	return id.String()
}

func WithRequestID(ctx context.Context) context.Context {
	requestID := generateRequestID()
	return context.WithValue(ctx, keys.KeyRequestID, requestID)
}

func (s *CacheRefresher) refresh(scrapeInterval time.Duration) {
	initialCtx := context.Background()

	requestIdCtx := WithRequestID(initialCtx)
	ctxLogger := zapctx.WithLogger(requestIdCtx, s.logger)

	tr := global.Tracer(domain.ServiceName)
	_, span := tr.Start(ctxLogger, "musicsnap/daemon/refresher.refresh", trace.WithNewRoot())
	defer span.End()

	s.userCache.Clean()

	//if err != nil {
	//	s.logger.Error("failed to refresh cache", zap.Error(err))
	//	return
	//}

	time.Sleep(scrapeInterval)
}
