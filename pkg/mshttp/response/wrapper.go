package response

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

const (
	httpCodeTag = "HTTP_CODE"
	rqUrlTag    = "RQ_URL"
)

type HttpResponseWrapper struct {
	logger *zap.Logger
}

func NewHttpResponseWrapper(logger *zap.Logger) *HttpResponseWrapper {
	return &HttpResponseWrapper{logger: logger}
}

func (r *HttpResponseWrapper) HandleSuccess(ctx *gin.Context, code int, msg string) {
	body := NewHttpResponse(code)
	r.handleSuccess(ctx, code, msg, body)
}

func (r *HttpResponseWrapper) HandleSuccessWithMessage(ctx *gin.Context, code int, msg string) {
	body := NewHttpResponseWithMessage(code, msg)
	r.handleSuccess(ctx, code, msg, body)
}

func (r *HttpResponseWrapper) handleSuccess(ctx *gin.Context, code int, msg string, body any) {
	r.logger.Debug(
		msg,
		zap.String(httpCodeTag, strconv.Itoa(code)),
		zap.String(rqUrlTag, ctx.Request.URL.String()))
	ctx.JSON(code, body)
}
