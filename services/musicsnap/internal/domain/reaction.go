package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

const (
	LikeReaction    = "like"
	DislikeReaction = "dislike"
)

// Reaction: Реакция на рецензию
type Reaction struct {
	ID int

	UserID uuid.UUID
	// ReviewID references Review(ID)
	ReviewID int

	Type      string // "like", "dislike" и т.д.
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r Reaction) Validate() error {
	if r.Type != LikeReaction && r.Type != DislikeReaction {
		return errors.New("reaction type must be 'like' or 'dislike'")
	}
	if r.ReviewID == 0 {
		return errors.New("review ID cannot be empty")
	}
	if r.UserID == uuid.Nil {
		return errors.New("user ID cannot be empty")
	}
	return nil
}

type ReactionCount struct {
	likes    int
	dislikes int
}
