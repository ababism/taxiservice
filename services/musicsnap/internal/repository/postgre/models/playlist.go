package models

// TODO DEPRECATED
//type PlaylistItemModel struct {
//	ID            uuid.UUID  `db:"id"`
//	PlaylistID    uuid.UUID  `db:"playlist_id"`
//	PieceID       *uuid.UUID `db:"piece_id"`       // nullable
//	DescriptionID *uuid.UUID `db:"description_id"` // nullable
//	Type          string     `db:"type"`
//	CreatedAt     time.Time  `db:"created_at"`
//	UpdatedAt     time.Time  `db:"updated_at"`
//}
//
//func (m *PlaylistItemModel) ToValidDomain() domain.PlaylistItem {
//	var pieceID, descID uuid.UUID
//	if m.PieceID != nil {
//		pieceID = *m.PieceID
//	}
//	if m.DescriptionID != nil {
//		descID = *m.DescriptionID
//	}
//
//	return domain.PlaylistItem{
//		ID:            m.ID,
//		PlaylistID:    m.PlaylistID,
//		PieceID:       pieceID,
//		DescriptionID: descID,
//		Type:          m.Type,
//		CreatedAt:     m.CreatedAt,
//		UpdatedAt:     m.UpdatedAt,
//	}
//}
//
//func ToPlaylistItemModel(p domain.PlaylistItem) PlaylistItemModel {
//	var pieceID, descID *uuid.UUID
//	if p.PieceID != uuid.Nil {
//		pieceID = &p.PieceID
//	}
//	if p.DescriptionID != uuid.Nil {
//		descID = &p.DescriptionID
//	}
//
//	return PlaylistItemModel{
//		ID:            p.ID,
//		PlaylistID:    p.PlaylistID,
//		PieceID:       pieceID,
//		DescriptionID: descID,
//		Type:          p.Type,
//		CreatedAt:     p.CreatedAt,
//		UpdatedAt:     p.UpdatedAt,
//	}
//}
