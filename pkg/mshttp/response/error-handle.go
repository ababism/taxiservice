package response

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func (r *HttpResponseWrapper) HandleError(ctx *gin.Context, code int, err error) {
	body := NewHttpResponse(code)
	r.handleError(ctx, code, err.Error(), body)
}

func (r *HttpResponseWrapper) HandleErrorWithMessage(ctx *gin.Context, code int, err error) {
	msg := err.Error()
	body := NewHttpResponseWithMessage(code, msg)
	r.handleError(ctx, code, msg, body)
}

func (r *HttpResponseWrapper) handleError(ctx *gin.Context, code int, msg string, body any) {
	r.logger.Error(
		msg,
		zap.String(httpCodeTag, strconv.Itoa(code)),
		zap.String(rqUrlTag, ctx.Request.URL.String()))
	ctx.AbortWithStatusJSON(code, body)
}
