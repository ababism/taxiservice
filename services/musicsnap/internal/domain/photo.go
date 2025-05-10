package domain

import (
	"github.com/google/uuid"
	"time"
)

// Photo: Фото
type Photo struct {
	ID        uuid.UUID
	PhotoURL  string
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PhotoParams struct {
	PhotoURL *string
	ID       *uuid.UUID
}
