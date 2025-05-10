package models

import (
	"github.com/google/uuid"
	"music-snap/services/musicsnap/internal/domain"
	"time"
)

type PhotoModel struct {
	ID       uuid.UUID `db:"id"`
	PhotoURL string    `db:"photo_url"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *PhotoModel) ToDomain() domain.Photo {
	return domain.Photo{
		ID:        m.ID,
		PhotoURL:  m.PhotoURL,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToEventPhotoModel(e domain.Photo) PhotoModel {
	return PhotoModel{
		ID:        e.ID,
		PhotoURL:  e.PhotoURL,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
