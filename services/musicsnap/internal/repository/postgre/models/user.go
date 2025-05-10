package models

import (
	"github.com/google/uuid"
	"music-snap/services/musicsnap/internal/domain"
	"time"
)

type UserModel struct {
	ID            uuid.UUID `db:"id"`
	Nickname      string    `db:"nickname"`
	AvatarURL     string    `db:"avatar_url"`
	BackgroundURL string    `db:"background_url"`
	Bio           string    `db:"bio"`

	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type RoleModel struct {
	Role      string    `db:"role"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ToRolesDomain(roles []RoleModel) domain.Roles {
	var rolesList []string
	for _, role := range roles {
		rolesList = append(rolesList, role.Role)
	}
	return domain.NewRoles(rolesList)
}

func (m *UserModel) ToDomain(roles []RoleModel) domain.User {
	r := ToRolesDomain(roles)

	p := domain.Profile{
		ID:            m.ID,
		Nickname:      m.Nickname,
		AvatarURL:     m.AvatarURL,
		BackgroundURL: m.BackgroundURL,
		Bio:           m.Bio,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
	return domain.User{
		Profile:      p,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Roles:        r,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (m *UserModel) ToDomainWithoutRoles() domain.User {
	p := domain.Profile{
		ID:            m.ID,
		Nickname:      m.Nickname,
		AvatarURL:     m.AvatarURL,
		BackgroundURL: m.BackgroundURL,
		Bio:           m.Bio,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
	return domain.User{
		Profile:      p,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
	}
}

func (m *UserModel) ToProfileDomain() domain.Profile {
	return domain.Profile{
		ID:            m.ID,
		Nickname:      m.Nickname,
		AvatarURL:     m.AvatarURL,
		BackgroundURL: m.BackgroundURL,
		Bio:           m.Bio,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func ToUserModel(u domain.User) UserModel {
	return UserModel{
		ID:            u.ID,
		Nickname:      u.Nickname,
		AvatarURL:     u.AvatarURL,
		BackgroundURL: u.BackgroundURL,
		Bio:           u.Bio,

		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
