package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"music-snap/services/musicsnap/internal/service"
	"strings"

	"music-snap/services/musicsnap/internal/config"
	"music-snap/services/musicsnap/internal/handler/http/musicsnap"
	oapigen "music-snap/services/musicsnap/internal/handler/http/musicsnap/oapi"
)

const (
	httpPrefix = "api"
	version    = "1"
)

type Handler struct {
	logger         *zap.Logger
	cfg            *config.Config
	coursesHandler *musicsnap.MusicsnapHandler
	//userServiceProvider service.MusicSnapService
}

// HandleError is a sample error handler function
func HandleError(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}

func InitHandler(
	router gin.IRouter,
	logger *zap.Logger,
	middlewares []oapigen.MiddlewareFunc,
	musicSnapService service.MusicSnapService,
) {
	msService := musicsnap.NewHandler(logger, musicSnapService)

	ginOpts := oapigen.GinServerOptions{
		BaseURL:      fmt.Sprintf("%s/%s", httpPrefix, getVersion()),
		Middlewares:  middlewares,
		ErrorHandler: HandleError,
	}
	oapigen.RegisterHandlersWithOptions(router, msService, ginOpts)
}

func getVersion() string {
	return fmt.Sprintf("v%s", strings.Split(version, ".")[0])
}
