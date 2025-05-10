package service

import (
	c "context"
	"fmt"
	"github.com/google/uuid"
	global "go.opentelemetry.io/otel"
	"music-snap/pkg/app"
	d "music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/service/ports"
	"net/http"
	"reflect"
)

func (s reactionSvc) spanName(funcName string) string {
	return fmt.Sprintf("%s/%s.%s.%s", "musicsnap", "service", reflect.TypeOf(s).Name(), funcName)
}

func NewReactionSvc(reaction ports.ReactionRepository) ports.ReactionService {
	return reactionSvc{r: reaction}
}

var _ ports.ReactionService = &reactionSvc{}

type reactionSvc struct {
	r ports.ReactionRepository
}

func (r reactionSvc) ListReactions(ctx c.Context, reviewID int, pagination d.IDPagination) ([]d.Reaction, d.IDPagination, error) {
	//TODO implement me
	panic("implement me")
}

func (r reactionSvc) CountReactions(ctx c.Context, reviewID uuid.UUID) (d.ReactionCount, error) {
	//TODO implement me
	panic("implement me")
}

func (s reactionSvc) UpdateReaction(ctx c.Context, actor d.Actor, reaction d.Reaction) (d.Reaction, error) {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("UpdateReaction"))
	defer span.End()
	ToSpan(&span, actor)

	if !actor.HasRole(d.AdminRole) && actor.ID != reaction.UserID {
		return d.Reaction{},
			app.NewError(http.StatusForbidden, "user can't update other persons reaction",
				"actor do not have admin role to update other persons reaction", nil)
	}

	reviewUpdated, err := s.r.Update(ctx, reaction)
	if err != nil {
		return d.Reaction{}, err
	}
	return reviewUpdated, nil
}

func (s reactionSvc) CreateReaction(ctx c.Context, actor d.Actor, reaction d.Reaction) (d.Reaction, error) {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("CreateReaction"))
	defer span.End()
	ToSpan(&span, actor)

	if !actor.HasRole(d.AdminRole) && actor.ID != reaction.UserID {
		return d.Reaction{},
			app.NewError(http.StatusForbidden, "user can't create other persons reaction",
				"actor do not have admin role to create other persons reaction", nil)
	}

	reviewCreated, err := s.r.Create(ctx, reaction)
	if err != nil {
		return d.Reaction{}, err
	}
	return reviewCreated, nil
}

func (s reactionSvc) GetByReview(ctx c.Context, actor d.Actor, reviewID int) (d.Reaction, error) {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("GetFromReview"))
	defer span.End()
	ToSpan(&span, actor)

	reaction, err := s.r.GetFromActor(ctx, actor.ID, reviewID)
	if err != nil {
		return d.Reaction{}, app.NewError(http.StatusNotFound, "reaction not found",
			"reaction with given id does not exist", err)
	}
	return reaction, nil
}

func (s reactionSvc) RemoveReaction(ctx c.Context, actor d.Actor, reactionID int) error {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("RemoveReaction"))
	defer span.End()
	ToSpan(&span, actor)

	if !actor.HasRole(d.AdminRole) {
		return app.NewError(http.StatusForbidden, "user can't remove other persons reaction",
			"actor do not have admin role to remove other persons reaction", nil)
	}

	err := s.r.Delete(ctx, reactionID)
	if err != nil {
		return app.NewError(http.StatusNotFound, "reaction not found",
			"reaction with given id does not exist", err)
	}
	return nil
}
