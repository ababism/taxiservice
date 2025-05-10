package oapi

import (
	"github.com/google/uuid"
	"music-snap/pkg/app"
	"music-snap/services/musicsnap/internal/domain"
	"net/http"
)

func (act Actor) ToValidDomain() (domain.Actor, error) {
	resActor := domain.Actor{}

	if act.Jwt == nil {
		return domain.Actor{}, app.NewError(http.StatusBadRequest, "jwt is required", "jwt is nil in actor", nil)
	}
	resActor.Jwt = *act.Jwt

	if act.Id == nil {
		//return domain.Actor{}, app.NewError(http.StatusBadRequest, "id is required", "id is nil in actor", nil)
		resActor.ID = uuid.Nil
	} else {
		resActor.ID = *act.Id
	}

	if act.Mail != nil {
		resActor.Mail = string(*act.Mail)
	}
	if act.Nickname != nil {
		resActor.Nickname = string(*act.Nickname)
	}

	if act.Roles != nil {
		resActor.AddRoles(*act.Roles)
	}

	return resActor, nil
}

func (p User) ToValidDomain() (domain.User, error) {
	profile := domain.Profile{}

	if p.Nickname == nil {
		return domain.User{}, app.NewError(http.StatusBadRequest, "nickname is required", "nickname is nil in user", nil)
	} else {
		profile.Nickname = *p.Nickname
	}
	//if p.Bio == nil {
	//	return domain.User{}, app.NewError(http.StatusBadRequest, "bio is required", "bio is nil in user", nil)
	//} else {
	//	profile.Bio = string(*p.Bio)
	//}
	if p.AvatarUrl != nil {
		profile.AvatarURL = string(*p.AvatarUrl)
	}
	if p.BackgroundUrl != nil {
		profile.BackgroundURL = string(*p.BackgroundUrl)
	}
	if p.Bio != nil {
		profile.Bio = *p.Bio
	}
	if p.Id != nil {
		profile.ID = *p.Id
	} else {
		profile.ID = uuid.Nil
	}

	user := domain.User{}
	user.Profile = profile

	if p.Email == nil {
		return domain.User{}, app.NewError(http.StatusBadRequest, "mail is required", "mail is nil in user", nil)
	} else {
		user.Email = string(*p.Email)
	}
	if p.Roles != nil {
		user.Roles = domain.NewRoles(*p.Roles)
	}

	return user, nil

}

func (b Profile) ToValidDomain() (domain.Profile, error) {
	profile := domain.Profile{}

	if b.Nickname == nil {
		return domain.Profile{}, app.NewError(http.StatusBadRequest, "nickname is required", "nickname is nil in user", nil)
	} else {
		profile.Nickname = *b.Nickname
	}
	if b.AvatarUrl != nil {
		profile.AvatarURL = string(*b.AvatarUrl)
	}
	if b.BackgroundUrl != nil {
		profile.BackgroundURL = string(*b.BackgroundUrl)
	}
	if b.Bio != nil {
		profile.Bio = *b.Bio
	}
	if b.Id != nil {
		profile.ID = *b.Id
	} else {
		profile.ID = uuid.Nil
	}

	return profile, nil
}

func ToUUIDPaginationDomain(l *int, lastUUID *UUID) domain.UUIDPagination {
	p := domain.UUIDPagination{}

	if l != nil {
		p.Limit = *l
	} else {
		p.Limit = 20
	}

	if lastUUID != nil {
		p.LastUUID = *lastUUID
	} else {
		p.LastUUID = uuid.Nil
	}

	return p
}

func ToIDPaginationDomain(l *int, lastID *int) domain.IDPagination {
	p := domain.IDPagination{}

	if l != nil {
		p.Limit = *l
	} else {
		p.Limit = 20
	}

	if lastID != nil {
		p.LastID = *lastID
	} else {
		p.LastID = 0
	}

	return p
}

func (sub Subscription) ToValidDomain() (domain.Subscription, error) {
	var subscription domain.Subscription

	if sub.SubscriberId == nil {
		return domain.Subscription{}, app.NewError(http.StatusBadRequest, "subscriberId is required", "subscriberId is nil in user", nil)
	} else {
		subscription.SubscriberID = *sub.SubscriberId
	}
	if sub.SubscribedToId == nil {
		return domain.Subscription{}, app.NewError(http.StatusBadRequest, "subscribedToId is required", "subscribedToId is nil in user", nil)
	} else {
		subscription.FollowedID = *sub.SubscribedToId
	}

	if sub.Id != nil {
		subscription.ID = *sub.Id
	} else {
		subscription.ID = int(0)
	}

	return subscription, nil
}

func (r Review) ToDomain() (domain.Review, error) {
	var review domain.Review
	// required
	if r.UserId == nil {
		return domain.Review{}, app.NewError(http.StatusBadRequest, "userId is required", "userId is nil in review", nil)
	} else {
		review.UserID = *r.UserId
	}
	if r.Rating == nil {
		return domain.Review{}, app.NewError(http.StatusBadRequest, "rating is required", "rating is nil in review", nil)
	} else {
		review.Rating = *r.Rating
	}
	if r.PieceId != nil {
		review.PieceID = *r.PieceId
	} else {
		return domain.Review{}, app.NewError(http.StatusBadRequest, "pieceId is required", "pieceId is nil in review", nil)
	}

	if r.Content != nil {
		review.Content = *r.Content
	}
	if r.PhotoUrl != nil {
		review.PhotoURL = *r.PhotoUrl
	}

	if r.Published != nil {
		review.Published = *r.Published
	} else {
		review.Published = true
	}

	if r.Moderated != nil {
		review.Moderated = *r.Moderated
	} else {
		review.Moderated = false
	}

	return review, nil
}

func (r GetReviewsListParams) ToDomain() (domain.ReviewFilter, error) {
	orderAsc := true
	orderByRating := false
	ofSubs := false

	var includeProfiles bool
	if r.IncludeProfiles != nil {
		includeProfiles = *r.IncludeProfiles
	} else {
		includeProfiles = false
	}
	filter := domain.ReviewFilter{
		UserID:          r.UserId,
		PieceID:         r.PieceId,
		Rating:          r.Rating,
		Moderated:       r.Moderated,
		Published:       r.Published,
		IncludeProfiles: includeProfiles,
		OrderByRating:   &orderByRating,
		OrderAsc:        orderAsc,
		OfSubscriptions: &ofSubs,
	}

	return filter, nil
}

func (r GetReviewsSubscriptionsParams) ToDomain() (domain.ReviewFilter, error) {
	orderAsc := true
	orderByRating := false
	ofSubs := true
	var includeProfiles bool
	if r.IncludeProfiles != nil {
		includeProfiles = *r.IncludeProfiles
	} else {
		includeProfiles = false
	}
	filter := domain.ReviewFilter{
		PieceID:         r.PieceId,
		Rating:          r.Rating,
		Moderated:       r.Moderated,
		Published:       r.Published,
		IncludeProfiles: includeProfiles,
		OrderByRating:   &orderByRating,
		OrderAsc:        orderAsc,
		OfSubscriptions: &ofSubs,
	}

	return filter, nil
}

func (r Reaction) ToDomain() (domain.Reaction, error) {
	var review domain.Reaction

	if r.UserId == nil {
		return domain.Reaction{}, app.NewError(http.StatusBadRequest, "author of reaction is required", "userId is nil in reaction", nil)
	} else {
		review.UserID = *r.UserId
	}

	if r.ReviewId == nil {
		return domain.Reaction{}, app.NewError(http.StatusBadRequest, "review target is required", "review_id is nil in reaction", nil)
	} else {
		review.ReviewID = *r.ReviewId
	}

	if r.Type == nil {
		return domain.Reaction{}, app.NewError(http.StatusBadRequest, "reaction type is required", "type is nil in reaction", nil)
	} else if *r.Type != domain.LikeReaction && *r.Type != domain.DislikeReaction {
		return domain.Reaction{}, app.NewError(http.StatusBadRequest, "reaction type is invalid", "type is not like or dislike in reaction", nil)
	} else {
		review.Type = (string)(*r.Type)
	}

	return review, nil
}

//func (f GetBannerParams) ToValidDomain() domain.BannerFilter {
//	return domain.BannerFilter{
//		Feature: f.FeatureId,
//		TagID:   f.TagId,
//		Limit:   f.Limit,
//		Offset:  f.Offset,
//	}
//}

//func (a Actor) ToDomainWithRoles(roles []string) domain.Actor {
//	roles = append(roles, a.Roles)
//	return domain.NewActor(a.ID, roles)
//}
//func ToTimeDomain(t *time.Time) time.Time {
//	if t == nil {
//		return time.Time{}
//	}
//	return *t
//}
//func (a Actor) ToValidDomain() domain.Actor {
//	s := []string{a.Roles}
//	return domain.NewActor(a.ID, s)
//}
