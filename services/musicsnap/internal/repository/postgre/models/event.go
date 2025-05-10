package models

import (
	"music-snap/services/musicsnap/internal/domain"
	"time"
)

type EventModel struct {
	ID   int       `db:"id"`
	Name string    `db:"name"`
	Date time.Time `db:"date"`
	//CoverURL   string    `db:"cover_url"`
	TicketLink string    `db:"ticket_link"`
	MapLink    string    `db:"map_link"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (m *EventModel) ToDomain() domain.Event {
	return domain.Event{
		ID:   m.ID,
		Name: m.Name,
		Date: m.Date,
		//CoverURL:   m.CoverURL,
		TicketLink:   m.TicketLink,
		LocationLink: m.MapLink,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func ToEventModel(e domain.Event) EventModel {
	return EventModel{
		ID:   e.ID,
		Name: e.Name,
		Date: e.Date,
		//CoverURL:   e.CoverURL,
		TicketLink: e.TicketLink,
		MapLink:    e.LocationLink,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}
