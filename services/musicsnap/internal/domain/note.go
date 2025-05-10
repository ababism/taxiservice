package domain

import "time"

// Note: Пользовательский контент
type Note struct {
	ID       int
	Text     string
	PhotoURL string

	SourceID   int
	SourceType string

	CreatedAt time.Time
	UpdatedAt time.Time
}
