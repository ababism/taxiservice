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

func (h MusicsnapHandler) PostUsers(c *gin.Context, params oapi.PostUsersParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PostUsers"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	var payload oapi.PostUsersJSONRequestBody
	h.bindRequestBody(c, &payload)

	userPayload, err := payload.User.ToValidDomain()
	if err != nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "invalid request body", "request body", err))
		return
	}

	if payload.Password == nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "password is required", "password is nil in user", nil))
		return
	}

	user, err := h.s.User.Create(ctx, actor, userPayload, *payload.Password)

	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	resp := oapi.ToUserResponse(user)
	c.JSON(http.StatusOK, resp)
}

// GetUsersProfiles is SearchProfiles request handler
func (h MusicsnapHandler) GetUsersProfiles(c *gin.Context, params oapi.GetUsersProfilesParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("GetUsersProfiles"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	pagination := oapi.ToUUIDPaginationDomain(params.Limit, params.LastUuid)

	profiles, pagination, err := h.s.User.GetProfilesList(ctx, actor, *params.NicknameQuery, pagination)

	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	type Response struct {
		Profiles   []oapi.Profile      `json:"profiles"`
		Pagination oapi.UUIDPagination `json:"pagination"`
	}

	response := Response{
		Profiles:   oapi.ToProfilesResponse(profiles),
		Pagination: oapi.ToUUIDPaginationResponse(pagination),
	}

	c.JSON(http.StatusOK, response)
}

func (h MusicsnapHandler) GetUsersUserId(c *gin.Context, userId oapi.UUID, params oapi.GetUsersUserIdParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("GetUsersUserId"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	user, err := h.s.User.Get(ctx, actor, userId)

	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	resp := oapi.ToUserResponse(user)
	c.JSON(http.StatusOK, resp)
}

func (h MusicsnapHandler) PutUsersUserId(c *gin.Context, userId oapi.UUID, params oapi.PutUsersUserIdParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PutUsersUserId"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	var payload oapi.PostUsersJSONRequestBody
	h.bindRequestBody(c, &payload)
	payload.User.Id = &userId

	userPayload, err := payload.User.ToValidDomain()
	if err != nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "invalid request body", "request body", err))
		return
	}

	if payload.Password == nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "password is required", "password is nil in user", nil))
		return
	}

	user, err := h.s.User.Update(ctx, actor, userPayload, *payload.Password)

	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	resp := oapi.ToUserResponse(user)
	c.JSON(http.StatusOK, resp)
}

func (h MusicsnapHandler) GetUsersUserIdProfile(c *gin.Context, userId oapi.UUID, params oapi.GetUsersUserIdProfileParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("GetUsersUserIdProfile"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	profile, err := h.s.User.GetProfile(ctx, actor, userId)

	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	resp := oapi.ToProfileResponse(profile)
	c.JSON(http.StatusOK, resp)
}

func (h MusicsnapHandler) PutUsersUserIdProfile(c *gin.Context, userId oapi.UUID, params oapi.PutUsersUserIdProfileParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PutUsersUserIdProfile"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	var payload oapi.PutUsersUserIdProfileJSONRequestBody
	h.bindRequestBody(c, &payload)
	payload.Id = &userId

	profilePayload, err := payload.ToValidDomain()
	if err != nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "invalid request body", "request body", err))
		return
	}

	profile, err := h.s.User.UpdateProfile(ctx, actor, profilePayload)

	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	resp := oapi.ToProfileResponse(profile)
	c.JSON(http.StatusOK, resp)
}
