package deleter

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

type DBCleaner struct {
	stop                   chan bool
	notificationRepository *ports.NotificationRepository
	logger                 *zap.Logger
}

func New(logger *zap.Logger, repository *ports.NotificationRepository) *DBCleaner {
	return &DBCleaner{
		logger:                 logger,
		notificationRepository: repository,
		stop:                   make(chan bool)}
}

func (s *DBCleaner) stopCallback(ctx context.Context) error {
	s.stop <- true
	return nil
}

func (s DBCleaner) StopFunc() func(context.Context) error {
	return s.stopCallback
}

func (s *DBCleaner) Start(scrapeInterval time.Duration) {
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

func (s *DBCleaner) refresh(scrapeInterval time.Duration) {
	initialCtx := context.Background()

	requestIdCtx := WithRequestID(initialCtx)
	ctxLogger := zapctx.WithLogger(requestIdCtx, s.logger)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctxLogger, "musicsnap/daemon/refresher.refresh", trace.WithNewRoot())
	defer span.End()

	err := s.notificationRepository.DeleteOutdated(ctx, time.Hour*24*365)
	if err != nil {
		s.logger.Error("failed to refresh cache", zap.Error(err))
		return
	}

	time.Sleep(scrapeInterval)
}
