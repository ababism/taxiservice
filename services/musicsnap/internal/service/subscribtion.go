package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	global "go.opentelemetry.io/otel"
	"music-snap/pkg/app"
	d "music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/service/ports"
	"net/http"
	"reflect"
)

func (s subscriptionSvc) spanName(funcName string) string {
	//return "musicsnap/service." + reflect.TypeOf(s).NameQuery() + "." + funcName
	return fmt.Sprintf("%s/%s.%s.%s", "musicsnap", "service", reflect.TypeOf(s).Name(), funcName)
}

func NewSubscriptionSvc(userRepository ports.UserRepository, cache ports.ProfileCache) ports.SubscriptionSvc {
	return subscriptionSvc{r: userRepository, c: cache, name: "subscription"}
}

var _ ports.SubscriptionSvc = &subscriptionSvc{}

type subscriptionSvc struct {
	r    ports.UserRepository
	c    ports.ProfileCache
	name string
}

func (s subscriptionSvc) Create(ctx context.Context, followingActor d.Actor, sub d.Subscription) (d.Subscription, error) {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("Create"))
	defer span.End()

	ToSpan(&span, followingActor)

	if followingActor.ID != sub.SubscriberID && !followingActor.HasRole(d.AdminRole) {
		return d.Subscription{}, app.NewError(http.StatusUnauthorized, "can't create subscription from other user",
			fmt.Sprintf("can't create sub between other user %s following %s, admin rights needed", followingActor.ID, sub.FollowedID), nil)
	}

	sub, err := s.r.CreateSub(ctx, sub)
	if err != nil {
		return d.Subscription{}, app.NewError(http.StatusBadRequest, "bad request, can't create subscription",
			fmt.Sprintf("can't create sub between %s following %s", followingActor.ID, sub.FollowedID), err)
	}

	return sub, nil
}

func (s subscriptionSvc) Update(ctx context.Context, followingActor d.Actor, sub d.Subscription) (d.Subscription, error) {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("Update"))
	defer span.End()

	ToSpan(&span, followingActor)

	if followingActor.ID != sub.SubscriberID && !followingActor.HasRole(d.AdminRole) {
		return d.Subscription{}, app.NewError(http.StatusUnauthorized, "can't update subscription from other user",
			fmt.Sprintf("can't update sub between other user %s following %s, admin rights needed", followingActor.ID, sub.FollowedID), nil)
	}

	prevSub, err := s.r.GetSub(ctx, sub.SubscriberID, sub.FollowedID)
	if err != nil {
		return d.Subscription{}, app.NewError(http.StatusBadRequest, "subscription not found",
			fmt.Sprintf("can't find sub between %s following %s", sub.SubscriberID, sub.FollowedID), err)
	}

	prevSub.NotificationFlag = sub.NotificationFlag

	sub, err = s.r.UpdateSub(ctx, prevSub)
	if err != nil {
		return d.Subscription{}, app.NewError(http.StatusBadRequest, "bad request, can't update subscription",
			fmt.Sprintf("can't update sub between %v", sub), err)
	}

	return sub, nil

}

func (s subscriptionSvc) Get(ctx context.Context, followingActor d.Actor, followedID uuid.UUID) (d.Subscription, error) {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("Get"))
	defer span.End()

	ToSpan(&span, followingActor)

	sub, err := s.r.GetSub(ctx, followingActor.ID, followedID)
	if err != nil {
		return d.Subscription{}, app.NewError(http.StatusBadRequest, "subscription not found",
			fmt.Sprintf("can't find sub between %s following %s", followingActor.ID, followedID), err)
	}

	return sub, nil
}

func (s subscriptionSvc) Delete(ctx context.Context, followingActor d.Actor, followedID uuid.UUID) error {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("Delete"))
	defer span.End()

	ToSpan(&span, followingActor)

	sub, err := s.r.GetSub(ctx, followingActor.ID, followedID)
	if err != nil {
		return app.NewError(http.StatusBadRequest, "subscription for deletion not found",
			fmt.Sprintf("can't find sub between %s following %s to delete", followingActor.ID, followedID), err)
	}

	if followingActor.ID != sub.SubscriberID && !followingActor.HasRole(d.AdminRole) {
		return app.NewError(http.StatusUnauthorized, "can't delete subscription from other user",
			fmt.Sprintf("can't delete sub between other user %s following %s, admin rights needed", followingActor.ID, sub.FollowedID), nil)
	}

	_, err = s.r.DeleteSub(ctx, sub)
	if err != nil {
		return app.NewError(http.StatusBadRequest, "can't delete subscription",
			fmt.Sprintf("can't delete subcription %v", sub), err)
	}

	return nil
}

func (s subscriptionSvc) Block(ctx context.Context, actor d.Actor, other uuid.UUID) error {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("Block"))
	defer span.End()

	ToSpan(&span, actor)

	subTo, errTo := s.r.GetSub(ctx, actor.ID, other)
	subFrom, errFrom := s.r.GetSub(ctx, other, actor.ID)

	if errTo != nil && errFrom != nil {
		return app.NewError(http.StatusBadRequest, "subscription not found",
			fmt.Sprintf("can't find any subs between %s and %s to block", actor.ID, other), errors.Join(errTo, errFrom))
	}

	if (actor.ID != subTo.SubscriberID || actor.ID != subFrom.FollowedID) && !actor.HasRole(d.AdminRole) {
		return app.NewError(http.StatusUnauthorized, "can't update subscription from other user",
			fmt.Sprintf("can't update sub between other user %s following %s, admin rights needed", actor, subFrom.FollowedID), nil)
	}

	_, err := s.r.DeleteSub(ctx, subTo)
	if err != nil {
		return app.NewError(http.StatusBadRequest, "can't delete subscription",
			fmt.Sprintf("can't delete subcription %v", subTo), err)
	}
	_, err = s.r.DeleteSub(ctx, subFrom)
	if err != nil {
		return app.NewError(http.StatusBadRequest, "can't delete subscription",
			fmt.Sprintf("can't delete subcription %v", subFrom), err)
	}

	return nil
}

func (s subscriptionSvc) GetSubscriptions(ctx context.Context, actor d.Actor, subscriberID uuid.UUID, pag d.IDPagination) ([]d.Subscription, d.IDPagination, error) {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("GetSubscriptions"))
	defer span.End()

	ToSpan(&span, actor)

	subs, _, err := s.r.ListSubscriptions(ctx, subscriberID, uuid.Nil, pag)
	if err != nil {
		return nil, pag, app.NewError(http.StatusBadRequest, "subscription not found",
			fmt.Sprintf("can't find any subscriptions of user %s", subscriberID), err)
	}

	return subs, pag, nil
}

func (s subscriptionSvc) GetSubscribers(ctx context.Context, actor d.Actor, followedID uuid.UUID, pag d.IDPagination) ([]d.Subscription, d.IDPagination, error) {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("GetSubscriptions"))
	defer span.End()

	ToSpan(&span, actor)

	subs, _, err := s.r.ListSubscriptions(ctx, uuid.Nil, followedID, pag)
	if err != nil {
		return nil, pag, app.NewError(http.StatusBadRequest, "subscription not found",
			fmt.Sprintf("can't find any followers of user %s", followedID), err)
	}

	return subs, pag, nil
}
