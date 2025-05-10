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

func (h MusicsnapHandler) PostUsersUserIdBlock(c *gin.Context, userId oapi.UUID, params oapi.PostUsersUserIdBlockParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PostUsersUserIdBlock"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	err = h.s.Subscription.Block(ctx, actor, userId)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, http.NoBody)
}

func (h MusicsnapHandler) GetUsersUserIdSubscribers(c *gin.Context, userId oapi.UUID, params oapi.GetUsersUserIdSubscribersParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("GetUsersUserIdSubscribers"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	pag := oapi.ToIDPaginationDomain(params.Limit, params.LastId)

	subscribers, pag, err := h.s.Subscription.GetSubscribers(ctx, actor, userId, pag)

	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	subsPart := oapi.ToSubsResponse(subscribers)
	pagPart := oapi.ToIDPaginationResponse(pag)

	type Response struct {
		Subscribers []oapi.Subscription `json:"subscribers"`
		Pagination  oapi.IDPagination   `json:"pagination"`
	}

	c.JSON(http.StatusOK, Response{
		Subscribers: subsPart,
		Pagination:  pagPart,
	})
}

func (h MusicsnapHandler) GetUsersUserIdSubscriptions(c *gin.Context, userId oapi.UUID, params oapi.GetUsersUserIdSubscriptionsParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("GetUsersUserIdSubscriptions"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	pag := oapi.ToIDPaginationDomain(params.Limit, params.LastId)

	subscriptions, pag, err := h.s.Subscription.GetSubscriptions(ctx, actor, userId, pag)

	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	subscriptionsPart := oapi.ToSubsResponse(subscriptions)
	pagPart := oapi.ToIDPaginationResponse(pag)

	type Response struct {
		Subscriptions []oapi.Subscription `json:"subscriptions"`
		Pagination    oapi.IDPagination   `json:"pagination"`
	}

	c.JSON(http.StatusOK, Response{
		Subscriptions: subscriptionsPart,
		Pagination:    pagPart,
	})
}

func (h MusicsnapHandler) PostSubscriptions(c *gin.Context, params oapi.PostSubscriptionsParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PostSubscriptions"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	var payload oapi.PostSubscriptionsJSONRequestBody
	h.bindRequestBody(c, &payload)

	subPayload, err := payload.ToValidDomain()
	if err != nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "invalid request body", "request body", err))
		return
	}

	subscription, err := h.s.Subscription.Create(ctx, actor, subPayload)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	resp := oapi.ToSubscriptionResponse(subscription)
	c.JSON(http.StatusOK, resp)
}

func (h MusicsnapHandler) DeleteSubscriptionsFollowedId(c *gin.Context, followedId oapi.UUID, params oapi.DeleteSubscriptionsFollowedIdParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("DeleteSubscriptionsFollowedId"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	err = h.s.Subscription.Delete(ctx, actor, followedId)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, http.NoBody)
}

func (h MusicsnapHandler) GetSubscriptionsFollowedId(c *gin.Context, followedId oapi.UUID, params oapi.GetSubscriptionsFollowedIdParams) {

	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("GetSubscriptionsFollowedId"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	sub, err := h.s.Subscription.Get(ctx, actor, followedId)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	resp := oapi.ToSubscriptionResponse(sub)
	c.JSON(http.StatusOK, resp)
}

func (h MusicsnapHandler) PutSubscriptionsFollowedId(c *gin.Context, followedId oapi.UUID, params oapi.PutSubscriptionsFollowedIdParams) {

	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PutSubscriptionsFollowedId"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	var payload oapi.PutSubscriptionsFollowedIdJSONRequestBody
	h.bindRequestBody(c, &payload)

	newSub := domain.Subscription{
		SubscriberID: actor.ID,
		FollowedID:   followedId,
	}

	if payload.NotificationFlag == nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "invalid request body notification flag is empty", "request body notification flag shouldn't be nil", err))
	}
	newSub.NotificationFlag = *payload.NotificationFlag

	subscription, err := h.s.Subscription.Update(ctx, actor, newSub)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	resp := oapi.ToSubscriptionResponse(subscription)
	c.JSON(http.StatusOK, resp)
}
