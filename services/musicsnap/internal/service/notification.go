package service

import (
	c "context"
	"fmt"
	"github.com/google/uuid"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/service/ports"
	"reflect"
)

func (s notificationSvc) spanName(funcName string) string {
	return fmt.Sprintf("%s/%s.%s.%s", "musicsnap", "service", reflect.TypeOf(s).Name(), funcName)
}

func NewNotificationService(notificationRepository ports.NotificationRepository) ports.NotificationSvc {
	return notificationSvc{notificationRepository: notificationRepository}
}

var _ ports.NotificationSvc = &notificationSvc{}

type notificationSvc struct {
	notificationRepository ports.NotificationRepository
	name                   string
}

func (s notificationSvc) GetNotifications(ctx c.Context, actor domain.Actor, pagination domain.UUIDPagination) ([]domain.Notification, domain.UUIDPagination, error) {
	//TODO implement me
	panic("implement me")
}

func (s notificationSvc) Notify(ctx c.Context, notification domain.Notification) error {
	//TODO implement me
	panic("implement me")
}

func (s notificationSvc) NotifyMany(ctx c.Context, notification []domain.Notification) error {
	//TODO implement me
	panic("implement me")
}

func (s notificationSvc) NotifyUsers(ctx c.Context, notification domain.Notification, userIDs []uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s notificationSvc) MarkAsRead(ctx c.Context, notificationID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
