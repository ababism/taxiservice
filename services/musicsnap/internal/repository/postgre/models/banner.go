package models

import (
	"encoding/json"
	"music-snap/services/musicsnap/internal/domain"
	"time"
)

//func (b *Banner) Value() (driver.Value, error) {
//	return json.Marshal(b.Content)
//}
//
//func (b *Banner) Scan(value interface{}) error {
//	musicsnap, ok := value.([]byte)
//	if !ok {
//		return errors.New("type assertion to []byte failed")
//	}
//
//	return json.Unmarshal(musicsnap, &b.Content)
//}

type Banner struct {
	//UUID      uuid.UUID
	ID int `db:"id"`
	//Content *map[string]interface{} `db:"content"  sql:"type:jsonb"`

	Content  json.RawMessage `db:"content"`
	IsActive bool            `db:"is_active"`

	Feature int   `db:"feature"`
	Tags    []int `db:"tags"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

//	func (b *Banner) ToDomain() domain.Banner {
//		res := domain.Banner{
//			ID:        b.ID,
//			IsActive:  b.IsActive,
//			Feature:   b.Feature,
//			Tags:      b.Tags,
//			UpdatedAt: b.UpdatedAt,
//			CreatedAt: b.CreatedAt,
//		}
//
//		err := json.Unmarshal(b.Content, &res.Content)
//		if err != nil {
//			log.Println(err)
//		}
//		return res
//	}
func (b *Banner) ToDomain() domain.Banner {
	return domain.Banner{
		ID:        b.ID,
		Content:   b.Content,
		IsActive:  b.IsActive,
		Feature:   b.Feature,
		Tags:      b.Tags,
		UpdatedAt: b.UpdatedAt,
		CreatedAt: b.CreatedAt,
	}
}

func ToBannerModel(b domain.Banner) Banner {
	return Banner{
		ID:        b.ID,
		Content:   b.Content,
		IsActive:  b.IsActive,
		Feature:   b.Feature,
		Tags:      b.Tags,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
