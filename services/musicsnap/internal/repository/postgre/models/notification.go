package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"music-snap/services/musicsnap/internal/domain"
	"time"
)

type NotificationModel struct {
	ID        uuid.UUID       `db:"id"`
	UserID    uuid.UUID       `db:"user_id"`
	Type      string          `db:"type"`
	Message   json.RawMessage `db:"message"`
	Read      bool            `db:"read"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

func (m *NotificationModel) ToDomain() (domain.Notification, error) {
	jsonObject := make(map[string]interface{})
	if err := json.Unmarshal(m.Message, &jsonObject); err != nil {
		return domain.Notification{}, err
	}
	return domain.Notification{
		ID:             m.ID,
		UserIDReceiver: m.UserID,
		Type:           m.Type,
		Message:        jsonObject,
		Read:           m.Read,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}, nil
}

func ToNotificationModel(n domain.Notification) (NotificationModel, error) {
	message, err := json.Marshal(n.Message)
	if err != nil {
		return NotificationModel{}, err
	}
	return NotificationModel{
		ID:        n.ID,
		UserID:    n.UserIDReceiver,
		Type:      n.Type,
		Message:   message,
		Read:      n.Read,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}, nil
}
