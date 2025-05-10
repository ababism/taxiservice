package models

import (
	"github.com/google/uuid"
	"music-snap/services/musicsnap/internal/domain"
	"time"
)

type SubscriptionModel struct {
	ID               int       `db:"sub_id"`
	SubscriberID     uuid.UUID `db:"subscriber_id"`
	FollowedID       uuid.UUID `db:"followed_id"`
	NotificationFlag bool      `db:"notification_flag"`
	//NotificationFlags []byte          `db:"notification_flags" sql:"type:jsonb"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *SubscriptionModel) ToLightDomain() domain.Subscription {

	return domain.Subscription{
		ID:               m.ID,
		SubscriberID:     m.SubscriberID,
		FollowedID:       m.FollowedID,
		NotificationFlag: m.NotificationFlag,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		//ProfileOfInterest: profile,
	}
}

func (m *SubscriptionModel) ToDomain(p domain.Profile) domain.Subscription {

	return domain.Subscription{
		ID:                m.ID,
		SubscriberID:      m.SubscriberID,
		FollowedID:        m.FollowedID,
		NotificationFlag:  m.NotificationFlag,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		ProfileOfInterest: p,
	}
}
func ToSubscriptionModel(s domain.Subscription) SubscriptionModel {

	return SubscriptionModel{
		ID:               s.ID,
		SubscriberID:     s.SubscriberID,
		FollowedID:       s.FollowedID,
		NotificationFlag: s.NotificationFlag,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}
