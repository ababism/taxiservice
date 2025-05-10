package musicsnap

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
	"music-snap/pkg/app"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/handler/http/models"
	"music-snap/services/musicsnap/internal/handler/http/musicsnap/oapi"
	"net/http"
)

func (h MusicsnapHandler) ReceiveActor(c context.Context, initialActor *oapi.Actor) (domain.Actor, error) {

	zapctx.Logger(c).With(zap.String("method", "ReceiveActor"))

	actor, err := initialActor.ToValidDomain()
	if err != nil {
		return domain.Actor{}, err
	}
	actor, err = h.s.Auth.EnrichActor(c, actor)
	if err != nil {
		zapctx.Logger(c).Debug("ReceiveActor failed to enrich actor", zap.Error(err), zap.Any("initial actor", initialActor), zap.Any("actor", actor))
		return domain.Actor{}, err
	}

	zapctx.Logger(c).With(zap.Any("actor", actor))
	return actor, nil
}

func parseRolesFromToken(c *gin.Context, token *string) ([]string, error) {
	if token == nil {
		return make([]string, 0), app.NewError(http.StatusUnauthorized, "token is required", "token is required", nil)
	}
	switch *token {
	case "user_token":
		return []string{domain.UserRole}, nil
	case "admin_token":
		return []string{domain.AdminRole}, nil
	default:
		return make([]string, 0), app.NewError(http.StatusUnauthorized, "invalid token", "invalid token", nil)
	}
}

func NewActorFromToken(c *gin.Context, token *string) (domain.Actor, error) {
	roles, err := parseRolesFromToken(c, token)
	if err != nil {
		return domain.Actor{}, err
	}
	return domain.NewActorFromRoles(roles), nil

}

func AbortWithBadResponse(c *gin.Context, logger *zap.Logger, statusCode int, err error) {
	logger.Debug(fmt.Sprintf("%s: %d %s", c.Request.URL, statusCode, app.GetLastMessage(err)))
	c.AbortWithStatusJSON(statusCode, models.Error{Message: app.GetLastMessage(err)})
}

func AbortWithErrorResponse(c *gin.Context, logger *zap.Logger, statusCode int, message string) {
	logger.Error(fmt.Sprintf("%s: %d %s", c.Request.URL, statusCode, message))
	c.AbortWithStatusJSON(statusCode, models.Error{Message: message})
}

func MapErrorToCode(err error) int {
	return app.GetCode(err)
}

func (h MusicsnapHandler) abortWithBadResponse(c *gin.Context, statusCode int, err error) {
	h.logger.Debug(fmt.Sprintf("%s: %d %s", c.Request.URL, statusCode, app.GetLastMessage(err)))
	c.AbortWithStatusJSON(statusCode, models.Error{Message: app.GetLastMessage(err)})
}

func (h MusicsnapHandler) abortWithAutoResponse(c *gin.Context, err error) {
	h.logger.Debug(fmt.Sprintf("%s: %d %s", c.Request.URL, app.GetCode(err), app.GetLastMessage(err)))
	c.AbortWithStatusJSON(app.GetCode(err), models.Error{Message: app.GetLastMessage(err)})
}

func (h MusicsnapHandler) bindRequestBody(c *gin.Context, obj any) bool {
	if err := c.BindJSON(obj); err != nil {
		AbortWithBadResponse(c, h.logger, http.StatusBadRequest, err)
		return false
	}
	return true
}
