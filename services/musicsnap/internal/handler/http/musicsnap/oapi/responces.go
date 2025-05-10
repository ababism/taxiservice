package oapi

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
	"music-snap/services/musicsnap/internal/domain"
)

func ToUserResponse(user domain.User) User {
	roles := user.Roles.ToSlice()
	email := (openapi_types.Email)(user.Email)
	return User{
		AvatarUrl:     &user.AvatarURL,
		BackgroundUrl: &user.BackgroundURL,
		Bio:           &user.Bio,
		CreatedAt:     &user.CreatedAt,
		Email:         &email,
		Id:            &user.ID,
		Nickname:      &user.Nickname,
		Roles:         &roles,
	}
}

func ToProfilesResponse(profiles []domain.Profile) []Profile {
	res := make([]Profile, len(profiles))
	for i, p := range profiles {
		res[i] = Profile{
			AvatarUrl:     &p.AvatarURL,
			BackgroundUrl: &p.BackgroundURL,
			Bio:           &p.Bio,
			CreatedAt:     &p.CreatedAt,
			Id:            &p.ID,
			Nickname:      &p.Nickname,
		}
	}
	return res
}

func ToProfileResponse(profile domain.Profile) Profile {
	return Profile{
		AvatarUrl:     &profile.AvatarURL,
		BackgroundUrl: &profile.BackgroundURL,
		Bio:           &profile.Bio,
		CreatedAt:     &profile.CreatedAt,
		Id:            &profile.ID,
		Nickname:      &profile.Nickname,
	}
}

func ToUUIDPaginationResponse(pag domain.UUIDPagination) UUIDPagination {
	return UUIDPagination{
		Limit:    &pag.Limit,
		LastUuid: &pag.LastUUID,
	}
}

func ToIDPaginationResponse(pag domain.IDPagination) IDPagination {
	return IDPagination{
		Limit:  &pag.Limit,
		LastId: &pag.LastID,
	}
}

func ToSubscriptionResponse(subs domain.Subscription) Subscription {
	pr := ToProfileResponse(subs.ProfileOfInterest)
	res := Subscription{
		CreatedAt:         &subs.CreatedAt,
		Id:                &subs.ID,
		ProfileOfInterest: &pr,
		SubscribedToId:    &subs.FollowedID,
		SubscriberId:      &subs.SubscriberID,
	}
	return res
}

func ToSubsResponse(subs []domain.Subscription) []Subscription {
	res := make([]Subscription, len(subs))
	for i, s := range subs {
		pr := ToProfileResponse(s.ProfileOfInterest)
		res[i] = Subscription{
			CreatedAt:         &s.CreatedAt,
			Id:                &s.ID,
			ProfileOfInterest: &pr,
		}
	}
	return res
}

func ToReviewResponse(review domain.Review) Review {
	var pr Profile
	if review.Profile != nil {
		pr = ToProfileResponse(*review.Profile)
	}
	return Review{
		Id: &review.ID,

		UserId:  &review.UserID,
		Profile: &pr,
		PieceId: &review.PieceID,

		Rating:   &review.Rating,
		Content:  &review.Content,
		PhotoUrl: &review.PhotoURL,

		Moderated: &review.Moderated,
		Published: &review.Published,

		CreatedAt: &review.CreatedAt,
		UpdatedAt: &review.UpdatedAt,
	}
}
func ToReviewsResponse(reviews []domain.Review) []Review {
	res := make([]Review, len(reviews))
	for i, r := range reviews {
		res[i] = ToReviewResponse(r)
	}
	return res
}

func ToReactionResponse(reaction domain.Reaction) Reaction {
	reactT := ReactionType((string)(reaction.Type))
	return Reaction{
		Id:        &reaction.ID,
		UserId:    &reaction.UserID,
		ReviewId:  &reaction.ReviewID,
		Type:      &reactT,
		CreatedAt: &reaction.CreatedAt,
	}
}

func ToReactionsResponse(reactions []domain.Reaction) []Reaction {
	res := make([]Reaction, len(reactions))
	for i, r := range reactions {
		res[i] = ToReactionResponse(r)
	}
	return res
}
