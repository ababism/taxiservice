package service

import (
	"music-snap/services/musicsnap/internal/repository/postgre"
	"music-snap/services/musicsnap/internal/service/ports"
)

// MusicSnapService is a struct that implements the MusicSnap service layer interfaces.
type MusicSnapService struct {
	Notification ports.NotificationSvc
	Auth         ports.AuthSvc
	User         ports.UserSvc
	Subscription ports.SubscriptionSvc
	Review       ports.ReviewService
	Reaction     ports.ReactionService
	Photo        ports.PhotoService
	Stats        ports.StatsService

	Event    ports.EventService
	Note     ports.NoteSvc
	Playlist ports.PlaylistService
}

func New(r postgre.Repository, jwt ports.JwtSvc, cache ports.ProfileCache) MusicSnapService {

	//notification := NewNotificationService(r.Notification)

	auth := NewAuthSvc(jwt, r.User)
	user := NewUserSvc(r.User, jwt, cache)
	subscription := NewSubscriptionSvc(r.User, cache)
	review := NewReviewSvc(r.Review, cache)
	reaction := NewReactionSvc(r.Reaction)
	// TODO
	//reaction := NewReactionSvc(r.Reaction)
	//photo := NewPhotoSvc(r.Photo)

	return MusicSnapService{
		//Notification: notification,

		Auth:         auth,
		User:         user,
		Subscription: subscription,

		Review:   review,
		Reaction: reaction,
		//Photo:    photo,

		//Event:  event,
		//Note:   note,
		//Banner: banner,
	}
}

func NewMSService(
	notification ports.NotificationSvc,
	auth ports.AuthSvc,
	user ports.UserSvc,
	subscription ports.SubscriptionSvc,
	review ports.ReviewService,
	reaction ports.ReactionService,
	photo ports.PhotoService,
	stats ports.StatsService,
	event ports.EventService,
	//note ports.NoteSvc,
	playlist ports.PlaylistService,
) *MusicSnapService {
	return &MusicSnapService{
		Notification: notification,
		Auth:         auth,
		User:         user,
		Subscription: subscription,
		Review:       review,
		Reaction:     reaction,
		Photo:        photo,
		Event:        event,
		//Note:         note,
	}
}

//type MusicSnapService struct {
//
//	ports.BannerService
//	ports.AuthSvc
//	ports.SubscriptionSvc
//	//ports.SongService
//	//ports.AlbumService
//	//ports.ArtistService
//	ports.PlaylistService
//	ports.ReviewService
//	ports.EventService
//	ports.NotificationSvc
//	ports.ReactionService
//	ports.CommentService
//	//ports.TagService
//	//ports.FollowService
//	ports.RatingService
//	//ports.StatisticsService
//	//ports.SearchService
//	//ports.PieceService
//}
