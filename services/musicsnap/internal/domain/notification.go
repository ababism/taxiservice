package domain

import (
	"github.com/google/uuid"
	"time"
)

// Notification: Уведомление
type Notification struct {
	ID uuid.UUID

	UserIDReceiver uuid.UUID
	UserIDSender   *uuid.UUID

	Type    string
	Message map[string]interface{}
	//Message json.RawMessage

	Read      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
