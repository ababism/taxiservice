package models

// TODO DEPRECATED
//type RatingModel struct {
//	ID        uuid.UUID `db:"id"`
//	UserID    uuid.UUID `db:"user_id"`
//	PieceID   uuid.UUID `db:"piece_id"`
//	Rating    int       `db:"rating"`
//	CreatedAt time.Time `db:"created_at"`
//	UpdatedAt time.Time `db:"updated_at"`
//}
//
//func (m *RatingModel) ToValidDomain() domain.Rating {
//	return domain.Rating{
//		ID:        m.ID,
//		UserID:    m.UserID,
//		PieceID:   m.PieceID,
//		Rating:    m.Rating,
//		CreatedAt: m.CreatedAt,
//		UpdatedAt: m.UpdatedAt,
//	}
//}
//
//func ToRatingModel(r domain.Rating) RatingModel {
//	return RatingModel{
//		ID:        r.ID,
//		UserID:    r.UserID,
//		PieceID:   r.PieceID,
//		Rating:    r.Rating,
//		CreatedAt: r.CreatedAt,
//		UpdatedAt: r.UpdatedAt,
//	}
//}
