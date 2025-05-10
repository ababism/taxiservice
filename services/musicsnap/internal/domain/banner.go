package domain

import (
	"encoding/json"
	"time"
)

type Banner struct {
	//UUID      uuid.UUID
	ID int

	json.RawMessage
	//Content *map[string]interface{}

	Content  json.RawMessage
	IsActive bool

	Feature int
	Tags    []int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CachedBanner struct {
	ID int

	//Content *map[string]interface{}
	Content json.RawMessage

	IsActive bool

	Expiration int64
	//CreatedAt time.Time
}

type BannerFilter struct {
	Feature *int
	TagID   *int

	Limit  *int
	Offset *int
}
