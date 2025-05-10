package models

//type CommentModel struct {
//	ID        uuid.UUID `db:"id"`
//	UserID    uuid.UUID `db:"user_id"`
//	ThreadID  uuid.UUID `db:"thread_id"`
//	Content      string    `db:"text"`
//	CreatedAt time.Time `db:"created_at"`
//	UpdatedAt time.Time `db:"updated_at"`
//}
//
//func (m *CommentModel) ToValidDomain() domain.Comment {
//	return domain.Comment{
//		ID:        m.ID,
//		UserID:    m.UserID,
//		ThreadID:  m.ThreadID,
//		Content:      m.Content,
//		CreatedAt: m.CreatedAt,
//		UpdatedAt: m.UpdatedAt,
//	}
//}
//
//func ToCommentModel(c domain.Comment) CommentModel {
//	return CommentModel{
//		ID:        c.ID,
//		UserID:    c.UserID,
//		ThreadID:  c.ThreadID,
//		Content:      c.Content,
//		CreatedAt: c.CreatedAt,
//		UpdatedAt: c.UpdatedAt,
//	}
//}
