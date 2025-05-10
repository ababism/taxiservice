package musicsnap

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/zaputil/zapctx"
	global "go.opentelemetry.io/otel"
	"music-snap/pkg/app"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/handler/http/musicsnap/oapi"
	"net/http"
)

func (h MusicsnapHandler) PostAuthLogin(c *gin.Context, params oapi.PostAuthLoginParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PostUsers"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	var payload oapi.PostAuthLoginJSONBody
	h.bindRequestBody(c, &payload)

	email := string(payload.Email)

	user, jwt, err := h.s.Auth.Login(ctx, actor, email, payload.Password)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	type Response struct {
		Jwt  string `json:"jwt"`
		User oapi.User
	}
	resp := Response{
		Jwt:  jwt,
		User: oapi.ToUserResponse(user),
	}
	c.JSON(http.StatusOK, resp)
}

func (h MusicsnapHandler) PostAuthLogout(c *gin.Context, params oapi.PostAuthLogoutParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PostUsers"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	_, _, err = h.s.Auth.LogOut(ctx, actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, http.NoBody)
}

func (h MusicsnapHandler) PostAuthRegister(c *gin.Context) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PostUsers"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	//h.logger.Debug("PostAuthRegister started")

	var payload oapi.PostAuthRegisterJSONRequestBody
	h.bindRequestBody(c, &payload)

	//h.logger.Debug("PostAuthRegister body parsed", zap.Any("payload", payload))

	userPayload, err := payload.User.ToValidDomain()
	if err != nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "invalid request body", "request body", err))
		return
	}

	if payload.Password == nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "password is required", "password is nil in user", nil))
		return
	}

	actor := domain.Actor{}
	jwt, user, err := h.s.Auth.Register(ctx, actor, userPayload, *payload.Password)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	type Response struct {
		Jwt  string `json:"jwt"`
		User oapi.User
	}
	resp := Response{
		Jwt:  jwt,
		User: oapi.ToUserResponse(user),
	}
	c.JSON(http.StatusOK, resp)
}
