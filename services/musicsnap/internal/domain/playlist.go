package domain

import (
	"github.com/google/uuid"
	"time"
)

// Playlist: Плейлист пользователя
type Playlist struct {
	ID          int
	UserID      uuid.UUID
	Name        string
	Description string

	// NotesIDs IDs для типа Note
	NotesIDs []int

	CoverURL string

	IsRanked  bool
	IsPrivate bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TODO DEPRECATED
//// PlaylistItem: Элемент плейлиста (Track, Album, Note)
//type PlaylistItem struct {
//	ID            uuid.UUID
//	PlaylistID    uuid.UUID
//	PieceID       uuid.UUID // Может быть nil, если Type = "Note"
//	DescriptionID uuid.UUID // Может быть nil, если Type = "Piece"
//	Type          string    // "Piece" или "Note"
//	CreatedAt     time.Time
//	UpdatedAt     time.Time
//}
