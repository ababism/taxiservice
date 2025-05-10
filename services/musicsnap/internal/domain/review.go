package domain

import (
	"github.com/google/uuid"
	"time"
)

// Review: Рецензия с описанием
type Review struct {
	ID int

	UserID  uuid.UUID
	Profile *Profile
	PieceID string

	Rating   int // 1-10
	Content  string
	PhotoURL string

	// For global
	Moderated bool
	Published bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ReviewFilter struct {
	UserID  *uuid.UUID
	PieceID *string

	Rating    *int // 1-10
	Moderated *bool
	Published *bool

	IncludeProfiles bool
	OrderByRating   *bool
	OrderAsc        bool

	OfSubscriptions *bool // true - only reviews of subscriptions
}

// TODO DEPRECATED
//// Thread: Тред комментариев
//type Thread struct {
//	ID        uuid.UUID
//	CreatedAt time.Time
//	UpdatedAt time.Time
//}

//// Comment: Комментарий в треде
//type Comment struct {
//	ID        uuid.UUID
//	UserID    uuid.UUID
//	ThreadID  uuid.UUID
//	Content      string
//	CreatedAt time.Time
//	UpdatedAt time.Time
//}
//
//// Rating: Оценка пользователя
//type Rating struct {
//	ID        uuid.UUID
//	UserID    uuid.UUID
//	PieceID   uuid.UUID
//	Rating    int // 1-10
//	CreatedAt time.Time
//	UpdatedAt time.Time
//}
