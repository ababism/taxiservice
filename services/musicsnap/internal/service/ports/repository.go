package ports

import (
	c "context"
	"github.com/google/uuid"
	d "music-snap/services/musicsnap/internal/domain"
	"time"
)

//type Repository struct {
//	User   UserRepository
//	Review ReviewRepository
//}

type BannerRepository interface {
	Find(ctx c.Context, tag, feature int) (d.Banner, error)

	Get(ctx c.Context, ID int) (d.Banner, error)
	Create(ctx c.Context, course d.Banner) (d.Banner, error)
	GetList(ctx c.Context, filter d.BannerFilter) ([]d.Banner, error)
	Update(ctx c.Context, bannerID int, banner d.Banner) (d.Banner, error)
	Delete(ctx c.Context, bannerID int) error
	ProcessDeleteQueue(ctx c.Context) error
}

// UserRepository: Управление пользователями
type UserRepository interface {
	Create(ctx c.Context, user d.User) (d.User, error)
	GetByID(ctx c.Context, id uuid.UUID) (d.User, error)
	GetByEmail(ctx c.Context, email string) (d.User, error)
	GetByNickname(ctx c.Context, nickname string) (d.User, error)
	GetByParam(ctx c.Context, filter d.UserFilter) (d.User, error)
	UpdateUser(ctx c.Context, user d.User) (d.User, error)
	Delete(ctx c.Context, id uuid.UUID) error
	GetList(ctx c.Context, nickNameQuery string, pag d.UUIDPagination) ([]d.User, uuid.UUID, error)
	GetProfile(ctx c.Context, profileID uuid.UUID) (d.Profile, error)

	CreateSub(ctx c.Context, sub d.Subscription) (d.Subscription, error)
	GetSub(ctx c.Context, subscriberID uuid.UUID, followedID uuid.UUID) (d.Subscription, error)
	UpdateSub(ctx c.Context, sub d.Subscription) (d.Subscription, error)
	DeleteSub(ctx c.Context, sub d.Subscription) (d.Subscription, error)
	ListSubscriptions(ctx c.Context, subscriberID uuid.UUID, followedID uuid.UUID, pag d.IDPagination) ([]d.Subscription, d.IDPagination, error)
}

// ReviewRepository: Управление рецензиями
type ReviewRepository interface {
	Create(ctx c.Context, review d.Review) (d.Review, error)
	Update(ctx c.Context, review d.Review) (d.Review, error)
	GetByID(ctx c.Context, id int) (d.Review, error)
	GetList(ctx c.Context, filter d.ReviewFilter, pag d.IDPagination) ([]d.Review, d.IDPagination, error)
	Delete(ctx c.Context, id int) (d.Review, error)
	//CreateReaction(ctx c.Context, reaction d.Reaction) error
	//GetComments(ctx c.Context, threadID uuid.UUID) ([]d.Comment, error)
}

// ReactionRepository: Управление реакциями
type ReactionRepository interface {
	Create(ctx c.Context, reaction d.Reaction) (d.Reaction, error)
	GetByID(ctx c.Context, id int) (d.Reaction, error)
	GetFromActor(ctx c.Context, userID uuid.UUID, reviewID int) (d.Reaction, error)
	GetFromReview(ctx c.Context, reviewID int) ([]d.Reaction, error)
	Update(ctx c.Context, reaction d.Reaction) (d.Reaction, error)
	Delete(ctx c.Context, id int) error
}

// PhotoRepository: Управление фотографиями событий
type PhotoRepository interface {
	Create(ctx c.Context, photo d.Photo) (d.Photo, error)
	GetByID(ctx c.Context, id uuid.UUID) (*d.Photo, error)
	GetByEvent(ctx c.Context, eventID uuid.UUID) ([]d.Photo, error)
	Update(ctx c.Context, photo d.Photo) error
	Delete(ctx c.Context, id uuid.UUID) error
}

type ProfileCache interface {
	Get(tag, feature int) (*d.CachedBanner, bool)
	Delete(tag, feature int)
	Set(tag, feature int, banner d.CachedBanner, duration time.Duration)
	Clean()
}

// EventRepository: Управление событиями
type EventRepository interface {
	Create(ctx c.Context, event d.Event) error
	GetByID(ctx c.Context, id int) (d.Event, error)
	Update(ctx c.Context, event d.Event) error
	Delete(ctx c.Context, id int) error
	ListByAuthor(ctx c.Context, authorID uuid.UUID) ([]d.Event, error)
	AddPhoto(ctx c.Context, photo d.Photo) error
}

// PlaylistRepository: Управление плейлистами
type PlaylistRepository interface {
	Create(ctx c.Context, playlist d.Playlist) error
	GetByID(ctx c.Context, id int) (d.Playlist, error)
	ListByUser(ctx c.Context, userID uuid.UUID, isPrivate bool) ([]d.Playlist, error)
	RemoveItem(ctx c.Context, playlistID, itemID uuid.UUID) error
}

// NotificationRepository: Управление уведомлениями
type NotificationRepository interface {
	Create(ctx c.Context, notification d.Notification) error
	MarkAsRead(ctx c.Context, id uuid.UUID) error
	ListByUser(ctx c.Context, userID uuid.UUID, unreadOnly bool) ([]d.Notification, error)
}

// DescriptionRepository: Управление описаниями
type DescriptionRepository interface {
	Create(ctx c.Context, description d.Note) error
	GetByID(ctx c.Context, id uuid.UUID) (d.Note, error)
	Update(ctx c.Context, description d.Note) error
	Delete(ctx c.Context, id uuid.UUID) error
}

// DEPRECATED ------

//// SubscriptionRepository: Управление подписками
//type SubscriptionRepository interface {
//	Create(ctx c.Context, subscription d.Subscription) (d.Subscription, error)
//	GetByID(ctx c.Context, id uuid.UUID) (d.Subscription, error)
//	GetBySubscriber(ctx c.Context, subscriberID uuid.UUID) ([]d.Subscription, error)
//	GetBySubscribedTo(ctx c.Context, subscribedToID uuid.UUID) ([]d.Subscription, error)
//	Update(ctx c.Context, subscription d.Subscription) error
//	Delete(ctx c.Context, id uuid.UUID) error
//}
//// PlaylistItemRepository: Управление элементами плейлиста
//type PlaylistItemRepository interface {
//	Create(ctx c.Context, item d.PlaylistItem) error
//	GetByID(ctx c.Context, id uuid.UUID) (d.PlaylistItem, error)
//	GetByPlaylist(ctx c.Context, playlistID uuid.UUID) ([]d.PlaylistItem, error)
//	Update(ctx c.Context, item d.PlaylistItem) error
//	Delete(ctx c.Context, id uuid.UUID) error
//}
//// CommentRepository: Управление комментариями
//type CommentRepository interface {
//	Create(ctx c.Context, comment d.Comment) error
//	GetByID(ctx c.Context, id uuid.UUID) (d.Comment, error)
//	GetByThread(ctx c.Context, threadID uuid.UUID) ([]d.Comment, error)
//	Update(ctx c.Context, comment d.Comment) error
//	Delete(ctx c.Context, id uuid.UUID) error
//}

//// ThreadRepository: Управление тредами
//type ThreadRepository interface {
//	Create(ctx c.Context, thread d.Thread) error
//	GetByID(ctx c.Context, id uuid.UUID) (d.Thread, error)
//	Update(ctx c.Context, thread d.Thread) error
//	Delete(ctx c.Context, id uuid.UUID) error
//}

//// RatingRepository: Управление оценками
//type RatingRepository interface {
//	Create(ctx c.Context, rating d.Rating) error
//	GetByID(ctx c.Context, id uuid.UUID) (d.Rating, error)
//	GetByUser(ctx c.Context, userID uuid.UUID) ([]d.Rating, error)
//	GetByPiece(ctx c.Context, pieceID uuid.UUID) ([]d.Rating, error)
//	Update(ctx c.Context, rating d.Rating) error
//	Delete(ctx c.Context, id uuid.UUID) error
//}
