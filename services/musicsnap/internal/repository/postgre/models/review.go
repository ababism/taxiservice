package models

import (
	"github.com/google/uuid"
	"music-snap/services/musicsnap/internal/domain"
	"time"
)

type ReviewModel struct {
	ID int `db:"id"`

	UserID  uuid.UUID `db:"user_id"`
	PieceID string    `db:"piece_id"`

	Rating   int    `db:"rating"`
	PhotoURL string `db:"photo_url"`
	Content  string `db:"content"`

	// For global
	Moderated bool `db:"moderated"`
	Published bool `db:"published"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *ReviewModel) ToDomain() domain.Review {
	return domain.Review{
		ID:        m.ID,
		UserID:    m.UserID,
		PieceID:   m.PieceID,
		Rating:    m.Rating,
		Content:   m.Content,
		PhotoURL:  m.PhotoURL,
		Moderated: m.Moderated,
		Published: m.Published,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m *ReviewModel) ToDomainProfile(p domain.Profile) domain.Review {
	p.ID = m.UserID
	return domain.Review{
		ID:        m.ID,
		UserID:    m.UserID,
		PieceID:   m.PieceID,
		Rating:    m.Rating,
		Content:   m.Content,
		PhotoURL:  m.PhotoURL,
		Moderated: m.Moderated,
		Published: m.Published,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Profile:   &p,
	}
}

func ToReviewModel(r domain.Review) ReviewModel {
	return ReviewModel{
		ID:        r.ID,
		UserID:    r.UserID,
		PieceID:   r.PieceID,
		Rating:    r.Rating,
		PhotoURL:  r.PhotoURL,
		Content:   r.Content,
		Moderated: r.Moderated,
		Published: r.Published,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
