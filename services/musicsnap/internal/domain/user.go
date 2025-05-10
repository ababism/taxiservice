package domain

import (
	"github.com/google/uuid"
	"time"
)

// Profile: Профиль пользователя
type Profile struct {
	ID            uuid.UUID
	Nickname      string
	AvatarURL     string
	BackgroundURL string
	Bio           string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// User: Основная сущность пользователя
type User struct {
	// общие данные
	Profile
	// полные данные для самого пользователя
	Email        string
	PasswordHash string
	Roles        Roles
	//Password     string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserFilter struct {
	ID       *uuid.UUID
	Nickname *string
	Email    *string

	CreatedAt *time.Time
	UpdatedAt *time.Time

	Limit    *int
	LastUUID *uuid.UUID
}

// User Copy method
func (u *User) Copy() User {
	return User{
		Profile: Profile{
			ID:            u.ID,
			Nickname:      u.Nickname,
			AvatarURL:     u.AvatarURL,
			BackgroundURL: u.BackgroundURL,
			Bio:           u.Bio,
			CreatedAt:     u.CreatedAt,
			UpdatedAt:     u.UpdatedAt,
		},
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Roles:        u.Roles,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func (u *User) Valid() bool {
	//if u.ID == uuid.Nil {
	//	return false
	//}
	if u.Nickname == "" {
		return false
	}
	if u.Email == "" {
		return false
	}
	if u.PasswordHash == "" {
		return false
	}
	if u.Roles.Empty() {
		return false
	}

	return true
}

// Subscription: Подписки пользователей
type Subscription struct {
	ID           int
	SubscriberID uuid.UUID // Ссылка на User.ID
	FollowedID   uuid.UUID // Ссылка на User.ID
	//NotificationFlags map[string]interface{} // JSON-флаги
	NotificationFlag  bool // JSON-флаги
	ProfileOfInterest Profile
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
