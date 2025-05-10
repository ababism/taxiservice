package domain

import (
	"github.com/google/uuid"
	"time"
)

// Event: Событие (концерт)
type Event struct {
	ID   int
	Name string
	Date time.Time
	//CoverURL   string

	Text     string
	PhotoURL string

	PhotosURLs   []string
	TicketLink   string
	LocationLink string
	Location     string

	Participants []uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}

// EventFilter: Фильтр для поиска событий
type EventFilter struct {
	ID        int
	NameQuery string

	DateLeftBound  time.Time
	DateRightBound time.Time

	SortByCreatedAt bool
	SortByAmount    bool
}
