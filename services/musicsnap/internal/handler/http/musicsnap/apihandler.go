package musicsnap

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/zaputil/zapctx"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/handler/http/musicsnap/oapi"
	"music-snap/services/musicsnap/internal/service"
	//generated "music-snap/services/musicsnap/internal/handler/http/musicsnap/oapi"
)

func (h MusicsnapHandler) spanName(funcName string) string {
	return "musicsnap/handler." + funcName
}

func NewHandler(logger *zap.Logger, msService service.MusicSnapService) *MusicsnapHandler {
	return &MusicsnapHandler{logger: logger, s: msService}
}

var _ oapi.ServerInterface = &MusicsnapHandler{}

type MusicsnapHandler struct {
	logger *zap.Logger
	s      service.MusicSnapService
}

func (h MusicsnapHandler) GetUsersUserIdStats(c *gin.Context, userId oapi.UUID, params oapi.GetUsersUserIdStatsParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) GetEvents(c *gin.Context, params oapi.GetEventsParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) PostEvents(c *gin.Context, params oapi.PostEventsParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) GetEventsEventId(c *gin.Context, eventId int, params oapi.GetEventsEventIdParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) PutEventsEventId(c *gin.Context, eventId int, params oapi.PutEventsEventIdParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) PostEventsEventIdParticipate(c *gin.Context, eventId int, params oapi.PostEventsEventIdParticipateParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) PostNotes(c *gin.Context, params oapi.PostNotesParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) DeleteNotesNoteId(c *gin.Context, noteId int, params oapi.DeleteNotesNoteIdParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) GetNotesNoteId(c *gin.Context, noteId int, params oapi.GetNotesNoteIdParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) PutNotesNoteId(c *gin.Context, noteId int, params oapi.PutNotesNoteIdParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) PostPhotos(c *gin.Context, params oapi.PostPhotosParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) DeletePhotosPhotoId(c *gin.Context, photoId oapi.UUID, params oapi.DeletePhotosPhotoIdParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) GetPhotosPhotoId(c *gin.Context, photoId oapi.UUID, params oapi.GetPhotosPhotoIdParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) GetPlaylistsPlaylistIdNotes(c *gin.Context, playlistId int, params oapi.GetPlaylistsPlaylistIdNotesParams) {
	//TODO implement me
	panic("implement me")
}

func (h MusicsnapHandler) GetTracksTrackIdStats(c *gin.Context, trackId string, params oapi.GetTracksTrackIdStatsParams) {
	//TODO implement me
	panic("implement me")
}

// NOTES TODO
func (h MusicsnapHandler) GetReviewsReviewIdNotes(c *gin.Context, reviewId int, params oapi.GetReviewsReviewIdNotesParams) {
	//TODO implement me
	panic("implement me")
}

// STAT Count
func (h MusicsnapHandler) GetReviewsReviewIdReactions(c *gin.Context, reviewId int, params oapi.GetReviewsReviewIdReactionsParams) {

	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(c, h.spanName("GetReviewsReviewIdReactions"))
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	actor, err := h.ReceiveActor(ctx, params.Actor)
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	_, _ = h.s.Review.GetReview(ctx, actor, reviewId, "")
	if err != nil {
		h.abortWithAutoResponse(c, err)
		return
	}

	// TODO implement stats
	panic("not implemented yet")

	//resp := oapi.ToReactionResponse(_)
	//c.JSON(http.StatusOK, resp)
}
