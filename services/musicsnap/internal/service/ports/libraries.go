package ports

import (
	d "music-snap/services/musicsnap/internal/domain"
)

// JwtSvc: Сервис для работы с JWT токенами
type JwtSvc interface {
	Generate(user d.User) (string, error)
	Parse(token string) (d.Actor, error)
}
