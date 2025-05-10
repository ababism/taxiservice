package models

import (
	"github.com/google/uuid"
	"music-snap/services/musicsnap/internal/domain"
	"time"
)

type ReactionModel struct {
	ID       int       `db:"id"`
	UserID   uuid.UUID `db:"user_id"`
	ReviewID int       `db:"review_id"`

	Type      string    `db:"type"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *ReactionModel) ToDomain() domain.Reaction {
	return domain.Reaction{
		ID:        m.ID,
		UserID:    m.UserID,
		ReviewID:  m.ReviewID,
		Type:      m.Type,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToReactionModel(r domain.Reaction) ReactionModel {
	return ReactionModel{
		ID:        r.ID,
		UserID:    r.UserID,
		ReviewID:  r.ReviewID,
		Type:      r.Type,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
