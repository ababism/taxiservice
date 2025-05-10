package ports

import (
	c "context"
	"github.com/google/uuid"
	d "music-snap/services/musicsnap/internal/domain"
)

// NotificationSvc: Бизнес-логика уведомлений
type NotificationSvc interface {
	// No api endpoint
	Notify(ctx c.Context, notification d.Notification) error
	// No api endpoint
	NotifyMany(ctx c.Context, notification []d.Notification) error
	// No api endpoint
	NotifyUsers(ctx c.Context, notification d.Notification, userIDs []uuid.UUID) error

	// endpoint with pagination by UUID
	GetNotifications(ctx c.Context, actor d.Actor, pagination d.UUIDPagination) ([]d.Notification, d.UUIDPagination, error)
	MarkAsRead(ctx c.Context, notificationID uuid.UUID) error
}

// AuthSvc: Бизнес-логика аутентификации
type AuthSvc interface {
	Register(ctx c.Context, actor d.Actor, user d.User, pass string) (jwt string, created d.User, err error)
	// password in the body of request
	Login(ctx c.Context, actor d.Actor, email, password string) (d.User, string, error)
	// No api endpoint
	EnrichActor(ctx c.Context, actor d.Actor) (d.Actor, error)

	LogOut(ctx c.Context, actor d.Actor) (d.User, string, error) // JWT токен
}

// UserSvc: Бизнес-логика пользователей
type UserSvc interface {
	// Create - создание пользователя администратором без получения токена
	Create(ctx c.Context, actor d.Actor, user d.User, pass string) (d.User, error)
	// id in query
	Get(ctx c.Context, actor d.Actor, userID uuid.UUID) (d.User, error)
	Update(ctx c.Context, actor d.Actor, user d.User, pass string) (d.User, error)

	// id  in query
	GetProfile(ctx c.Context, actor d.Actor, userID uuid.UUID) (d.Profile, error)
	UpdateProfile(ctx c.Context, actor d.Actor, user d.Profile) (d.Profile, error)
	// api endpoint with pagination by UUID for profile search
	GetProfilesList(ctx c.Context, actor d.Actor,
		nickNameQuery string, pagination d.UUIDPagination) ([]d.Profile, d.UUIDPagination, error)
}

// SubscriptionSvc: Бизнес-логика подписок
type SubscriptionSvc interface {
	Create(ctx c.Context, followingActor d.Actor, sub d.Subscription) (d.Subscription, error)
	Update(ctx c.Context, followingActor d.Actor, sub d.Subscription) (d.Subscription, error)
	// followedID in query
	Get(ctx c.Context, followingActor d.Actor, followedID uuid.UUID) (d.Subscription, error)
	Delete(ctx c.Context, followingActor d.Actor, followedID uuid.UUID) error

	// other id in query
	Block(ctx c.Context, actor d.Actor, other uuid.UUID) error

	GetSubscriptions(ctx c.Context, actor d.Actor, subscriberID uuid.UUID, pag d.IDPagination) ([]d.Subscription, d.IDPagination, error)
	GetSubscribers(ctx c.Context, actor d.Actor, followedID uuid.UUID, pagination d.IDPagination) ([]d.Subscription, d.IDPagination, error)
}

// ReviewService: Бизнес-логика рецензий
type ReviewService interface {
	CreateReview(ctx c.Context, actor d.Actor, review d.Review) (d.Review, error)
	UpdateReview(ctx c.Context, actor d.Actor, review d.Review) (d.Review, error)
	GetReview(ctx c.Context, actor d.Actor, reviewID int, pieceID string) (d.Review, error)
	DeleteReview(ctx c.Context, actor d.Actor, reviewID int) error
	//ReviewsOfSubscriptions(ctx c.Context, actor d.Actor, filter d.ReviewFilter, pagination d.IDPagination) ([]d.Review, d.IDPagination, error)
	ListReviews(ctx c.Context, actor d.Actor, filter d.ReviewFilter, pagination d.IDPagination) ([]d.Review, d.IDPagination, error)
}

type ReactionService interface {
	CreateReaction(ctx c.Context, actor d.Actor, reaction d.Reaction) (d.Reaction, error)
	UpdateReaction(ctx c.Context, actor d.Actor, reaction d.Reaction) (d.Reaction, error)

	GetByReview(ctx c.Context, actor d.Actor, reviewID int) (d.Reaction, error)
	RemoveReaction(ctx c.Context, actor d.Actor, reactionID int) error
	ListReactions(ctx c.Context, reviewID int, pagination d.IDPagination) ([]d.Reaction, d.IDPagination, error)

	CountReactions(ctx c.Context, reviewID uuid.UUID) (d.ReactionCount, error)
}

// PhotoService: Бизнес-логика фотографий событий
type PhotoService interface {
	AddPhoto(ctx c.Context, actor d.Actor, photo d.Photo, file []byte) (d.Photo, error)
	RemovePhoto(ctx c.Context, actor d.Actor, photo d.Photo) error
	GetPhoto(ctx c.Context, actor d.Actor, photo d.PhotoParams) (d.Photo, []byte, error)
}

// StatsService: Бизнес-логика статистики
type StatsService interface {
	GetProfileStats(ctx c.Context, actor d.Actor, userID uuid.UUID) (d.ProfileStats, error)
	GetMusicTrackStats(ctx c.Context, actor d.Actor, trackID string) (d.TrackStats, error)

	// UpdateDWH For DWH daemon, no api calls
	UpdateDWH(ctx c.Context) error
}

// EventService: Бизнес-логика событий
type EventService interface {
	CreateEvent(ctx c.Context, actor d.Actor, event d.Event) (d.Event, error)
	GetEvent(ctx c.Context, eventID uuid.UUID) (d.Event, error)
	UpdateEvent(ctx c.Context, actor d.Actor, event d.Event) (d.Event, error)
	ListEvents(ctx c.Context, actor d.Actor, filters d.EventFilter) ([]d.Event, error)
	Participate(ctx c.Context, actor d.Actor, eventID uuid.UUID) (d.Event, error)
}

// NoteSvc: Бизнес-логика описаний
type NoteSvc interface {
	CreateNote(ctx c.Context, actor d.Actor, description d.Note) error
	GetNote(ctx c.Context, descriptionID uuid.UUID) (d.Note, error)
	UpdateNote(ctx c.Context, actor d.Actor, description d.Note) error
	DeleteNote(ctx c.Context, actor d.Actor, description d.Note) error
	GetPlaylistNotes(ctx c.Context, playlistID int) ([]d.Note, error)
	GetReviewCommentNotes(ctx c.Context, reviewID int, pagination d.IDPagination) ([]d.Note, d.IDPagination, error)
}

// PlaylistService: Бизнес-логика плейлистов
type PlaylistService interface {
	CreatePlaylist(ctx c.Context, actor d.Actor, playlist d.Playlist) (d.Playlist, error)
	GetPlaylist(ctx c.Context, playlistID uuid.UUID) (d.Playlist, error)
	UpdatePlaylist(ctx c.Context, actor d.Actor, playlist d.Playlist) (d.Playlist, error)
	DeletePlaylist(ctx c.Context, actor d.Actor, playlist d.Playlist) (d.Playlist, error)
	ListPlaylists(ctx c.Context, actor d.Actor, userID uuid.UUID, paginationID int) ([]d.Playlist, d.IDPagination, error)
}

// DONE ---------------

// Not today
//type BannerService interface {
//	Find(ctx c.Context, actor d.Actor, tag, feature int, useLastRevision bool) (d.CachedBanner, error)
//
//	Create(ctx c.Context, actor d.Actor, course d.Banner) (int, error)
//	GetList(ctx c.Context, actor d.Actor, filter d.BannerFilter) ([]d.Banner, error)
//	Update(ctx c.Context, actor d.Actor, bannerID int, banner d.Banner) error
//	Delete(ctx c.Context, actor d.Actor, bannerID int) error
//	//Get(ctx c.Context, actor d.Actor, ID int) (d.Banner, error)
//}

//// ThreadService: Бизнес-логика тредов
//type ThreadService interface {
//	CreateThread(ctx c.Context, actor d.Actor, thread d.Thread) error
//	GetThread(ctx c.Context, threadID uuid.UUID) (d.Thread, error)
//	UpdateThread(ctx c.Context, actor d.Actor, thread d.Thread) error
//	DeleteThread(ctx c.Context, actor d.Actor, threadID uuid.UUID) error
//}

// TODO DEPRECATED

//// CommentService: Бизнес-логика комментариев
//type CommentService interface {
//	CreateComment(ctx c.Context, actor d.Actor, comment d.Comment) error
//	GetComment(ctx c.Context, commentID uuid.UUID) (d.Comment, error)
//	UpdateComment(ctx c.Context, actor d.Actor, comment d.Comment) error
//	DeleteComment(ctx c.Context, actor d.Actor, commentID uuid.UUID) error
//	ListComments(ctx c.Context, threadID uuid.UUID) ([]d.Comment, error)
//}
//
//// PlaylistItemService: Бизнес-логика элементов плейлиста
//type PlaylistItemService interface {
//	AddToPlaylist(ctx c.Context, actor d.Actor, item d.PlaylistItem) error
//	RemoveFromPlaylist(ctx c.Context, actor d.Actor, playlistID, itemID uuid.UUID) error
//	GetPlaylistItems(ctx c.Context, playlistID uuid.UUID) ([]d.PlaylistItem, error)
//	UpdatePlaylistItem(ctx c.Context, actor d.Actor, item d.PlaylistItem) error
//}

//// RatingService: Бизнес-логика оценок
//type RatingService interface {
//	CreateRating(ctx c.Context, actor d.Actor, rating d.Rating) error
//	GetRating(ctx c.Context, ratingID uuid.UUID) (d.Rating, error)
//	UpdateRating(ctx c.Context, actor d.Actor, rating d.Rating) error
//	DeleteRating(ctx c.Context, actor d.Actor, ratingID uuid.UUID) error
//	ListRatingsByUser(ctx c.Context, userID uuid.UUID) ([]d.Rating, error)
//	ListRatingsByPiece(ctx c.Context, pieceID uuid.UUID) ([]d.Rating, error)
//}
