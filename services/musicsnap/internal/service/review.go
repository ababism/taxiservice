package service

import (
	c "context"
	"fmt"
	"github.com/google/uuid"
	global "go.opentelemetry.io/otel"
	"music-snap/pkg/app"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/service/ports"
	"net/http"
	"reflect"
)

func (s reviewSvc) spanName(funcName string) string {
	//return "musicsnap/service." + reflect.TypeOf(s).NameQuery() + "." + funcName
	return fmt.Sprintf("%s/%s.%s.%s", "musicsnap", "service", reflect.TypeOf(s).Name(), funcName)
}

func NewReviewSvc(reviewRepository ports.ReviewRepository, cache ports.ProfileCache) ports.ReviewService {
	return reviewSvc{r: reviewRepository, c: cache}
}

var _ ports.ReviewService = &reviewSvc{}

type reviewSvc struct {
	r   ports.ReviewRepository
	c   ports.ProfileCache
	jwt ports.JwtSvc
}

func (s reviewSvc) validForCreation(r domain.Review) error {
	if r.UserID == uuid.Nil {
		return app.NewError(http.StatusBadRequest, "user id is required", "user id is nil in review", nil)
	}
	if r.Rating < 1 || r.Rating > 10 {
		return app.NewError(http.StatusBadRequest, "rating must be between 1 and 10", "invalid rating in review", nil)
	}
	return nil
}
func (s reviewSvc) CreateReview(ctx c.Context, actor domain.Actor, review domain.Review) (domain.Review, error) {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("CreateReview"))
	defer span.End()
	ToSpan(&span, actor)

	if !actor.HasRole(domain.AdminRole) && actor.ID != review.UserID {
		return domain.Review{},
			app.NewError(http.StatusForbidden, "user can't create other persons review",
				"actor do not have admin role to create other persons review", nil)
	}

	err := s.validForCreation(review)
	if err != nil {
		return domain.Review{},
			app.NewError(http.StatusForbidden, "invalid review",
				"invalid fields for review validation", err)
	}

	reviewCreated, err := s.r.Create(ctx, review)
	if err != nil {
		return domain.Review{}, err
	}
	return reviewCreated, nil
}

func (s reviewSvc) UpdateReview(ctx c.Context, actor domain.Actor, review domain.Review) (domain.Review, error) {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("UpdateReview"))
	defer span.End()
	ToSpan(&span, actor)

	if !actor.HasRole(domain.AdminRole) && actor.ID != review.UserID {
		return domain.Review{},
			app.NewError(http.StatusForbidden, "user can't update other persons review",
				"actor do not have admin role to update other persons review", nil)
	}

	err := s.validForCreation(review)
	if err != nil {
		return domain.Review{},
			app.NewError(http.StatusForbidden, "invalid review",
				"invalid fields for review validation", err)
	}

	reviewUpdated, err := s.r.Update(ctx, review)
	if err != nil {
		return domain.Review{}, err
	}
	return reviewUpdated, nil
}

func (s reviewSvc) GetReview(ctx c.Context, actor domain.Actor, revID int, pieceID string) (domain.Review, error) {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("GetReview"))
	defer span.End()

	review, err := s.r.GetByID(ctx, revID)
	if err != nil {
		return domain.Review{}, err
	}
	return review, nil
}

func (s reviewSvc) DeleteReview(ctx c.Context, actor domain.Actor, reviewID int) error {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("DeleteReview"))
	defer span.End()
	ToSpan(&span, actor)

	review, err := s.r.GetByID(ctx, reviewID)
	if err != nil {
		return app.NewError(http.StatusNotFound, "review not found",
			fmt.Sprintf("review with id %s not found", reviewID), err)
	}

	// Check if the actor is allowed to delete the review
	if !actor.HasRole(domain.AdminRole) && actor.ID != review.UserID {
		return app.NewError(http.StatusForbidden, "user can't delete other persons review",
			"actor do not have admin role to delete other persons review", nil)
	}

	_, err = s.r.Delete(ctx, review.ID)
	if err != nil {
		return err
	}
	return nil
}

//func (s reviewSvc) ReviewsOfSubscriptions(ctx c.Context, actor domain.Actor, filter domain.ReviewFilter, pagination domain.IDPagination) ([]domain.Review, domain.IDPagination, error) {
//	//TODO implement me
//	panic("implement me")
//}

func (s reviewSvc) ListReviews(ctx c.Context, actor domain.Actor, filter domain.ReviewFilter, pagination domain.IDPagination) ([]domain.Review, domain.IDPagination, error) {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("ListReviews"))
	defer span.End()
	ToSpan(&span, actor)

	reviews, pag, err := s.r.GetList(ctx, filter, pagination)
	if err != nil {
		return nil, domain.IDPagination{}, app.NewError(http.StatusInternalServerError, "error listing reviews",
			"error while listing reviews", err)
	}

	return reviews, pag, nil
}
