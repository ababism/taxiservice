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

// REACTIONS
func (h MusicsnapHandler) PostReviewsReviewIdReactions(c *gin.Context, reviewId int, params oapi.PostReviewsReviewIdReactionsParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PostReviewsReviewIdReactions"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	var payload oapi.PostReviewsReviewIdReactionsJSONRequestBody
	h.bindRequestBody(c, &payload)

	reactionPayload, err := payload.ToDomain()
	if err != nil {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "invalid request body", "request body", err))
		return
	}

	reviewReaction, err := h.s.Reaction.CreateReaction(ctx, actor, reactionPayload)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	resp := oapi.ToReactionResponse(reviewReaction)
	c.JSON(http.StatusOK, resp)
}
func (h MusicsnapHandler) PutReactionsReactionId(c *gin.Context, reactionId int, params oapi.PutReactionsReactionIdParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("PutReactionsReactionId"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	if params.Type == "" {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "invalid request params", "reaction type is required", nil))
		return
	}
	if params.Type != domain.DislikeReaction && params.Type != domain.LikeReaction {
		h.abortWithAutoResponse(c, app.NewError(http.StatusBadRequest, "invalid request params", "reaction type must be like or dislike", nil))
		return
	}

	reaction := domain.Reaction{ID: reactionId, Type: params.Type}

	reviewUpdated, err := h.s.Reaction.UpdateReaction(ctx, actor, reaction)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	resp := oapi.ToReactionResponse(reviewUpdated)
	c.JSON(http.StatusOK, resp)
}

func (h MusicsnapHandler) DeleteReactionsReactionId(c *gin.Context, reactionId int, params oapi.DeleteReactionsReactionIdParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("DeleteReactionsReactionId"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	err = h.s.Reaction.RemoveReaction(ctx, actor, reactionId)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h MusicsnapHandler) GetReviewsReviewIdReactionsMe(c *gin.Context, reviewId int, params oapi.GetReviewsReviewIdReactionsMeParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("GetReviewsReviewIdReactions"))
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

	reactions, err := h.s.Reaction.GetByReview(ctx, actor, review.ID)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	resp := oapi.ToReactionResponse(reactions)
	c.JSON(http.StatusOK, resp)
}
