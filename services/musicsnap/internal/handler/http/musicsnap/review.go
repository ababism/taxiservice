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

func (h MusicsnapHandler) PostReviews(c *gin.Context, params oapi.PostReviewsParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PostReviews"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	var payload oapi.PostReviewsJSONRequestBody
	h.bindRequestBody(c, &payload)

	reviewPayload, err := payload.ToDomain()
	if err != nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "invalid request body", "request body", err))
		return
	}

	review, err := h.s.Review.CreateReview(ctx, actor, reviewPayload)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	resp := oapi.ToReviewResponse(review)
	c.JSON(http.StatusOK, resp)
}

func (h MusicsnapHandler) GetReviewsList(c *gin.Context, params oapi.GetReviewsListParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("GetReviewsList"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	filter, err := params.ToDomain()

	pag := oapi.ToIDPaginationDomain(params.Limit, params.LastId)

	reviews, pag, err := h.s.Review.ListReviews(ctx, actor, filter, pag)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	reviewsPart := oapi.ToReviewsResponse(reviews)
	pagPart := oapi.ToIDPaginationResponse(pag)

	type Response struct {
		Reviews    []oapi.Review     `json:"reviews"`
		Pagination oapi.IDPagination `json:"pagination"`
	}

	c.JSON(http.StatusOK, Response{
		Reviews:    reviewsPart,
		Pagination: pagPart,
	})
}

func (h MusicsnapHandler) GetReviewsSubscriptions(c *gin.Context, params oapi.GetReviewsSubscriptionsParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("GetReviewsSubscriptions"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	filter, err := params.ToDomain()
	if err != nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "invalid request body", "request body", err))
		return
	}

	pag := oapi.ToIDPaginationDomain(params.Limit, params.LastId)

	reviews, pag, err := h.s.Review.ListReviews(ctx, actor, filter, pag)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	reviewsPart := oapi.ToReviewsResponse(reviews)
	pagPart := oapi.ToIDPaginationResponse(pag)

	type Response struct {
		Reviews    []oapi.Review     `json:"reviews"`
		Pagination oapi.IDPagination `json:"pagination"`
	}

	c.JSON(http.StatusOK, Response{
		Reviews:    reviewsPart,
		Pagination: pagPart,
	})
}

func (h MusicsnapHandler) DeleteReviewsReviewId(c *gin.Context, reviewId int, params oapi.DeleteReviewsReviewIdParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("DeleteReviewsReviewId"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	err = h.s.Review.DeleteReview(ctx, actor, reviewId)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h MusicsnapHandler) GetReviewsReviewId(c *gin.Context, reviewId int, params oapi.GetReviewsReviewIdParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("GetReviewsReviewId"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	review, err := h.s.Review.GetReview(ctx, actor, reviewId, "")
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	resp := oapi.ToReviewResponse(review)
	c.JSON(http.StatusOK, resp)
}

func (h MusicsnapHandler) PutReviewsReviewId(c *gin.Context, reviewId int, params oapi.PutReviewsReviewIdParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PutReviewsReviewId"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	var payload oapi.PutReviewsReviewIdJSONRequestBody
	h.bindRequestBody(c, &payload)

	reviewPayload, err := payload.ToDomain()
	if err != nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "invalid request body", "request body", err))
		return
	}

	reviewPayload.ID = reviewId

	reviewUpdated, err := h.s.Review.UpdateReview(ctx, actor, reviewPayload)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	resp := oapi.ToReviewResponse(reviewUpdated)
	c.JSON(http.StatusOK, resp)
}
