package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	mserror "music-snap/pkg/app"
	"net/http"
)

type xorErrorHandler func(ctx *gin.Context, code int, err error)

func (r *HttpResponseWrapper) HandleMSError(ctx *gin.Context, err error) {
	handleXorError(ctx, err, r.HandleError, r.HandleError)
}

func (r *HttpResponseWrapper) HandleMSErrorWithMessage(ctx *gin.Context, err error) {
	handleXorError(ctx, err, r.HandleErrorWithMessage, r.HandleError)
}

func handleXorError(ctx *gin.Context, err error, handler xorErrorHandler, defaultHandler xorErrorHandler) {
	switch {
	case errors.As(err, &mserror.Error{}):
		handler(ctx, http.StatusBadRequest, err)
	default:
		defaultHandler(ctx, http.StatusInternalServerError, err)
	}
}
