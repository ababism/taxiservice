package jwtservice

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"music-snap/pkg/app"
	d "music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/service/ports"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	SigningKey string `mapstructure:"signingkey"`
	TTLHours   int    `mapstructure:"ttl_hours"`
}

var _ ports.JwtSvc = &Svc{}

type Svc struct {
	signingKey string
	ttl        time.Duration
}

func (c *Config) Validate() error {
	if c == nil {
		return app.NewError(http.StatusInternalServerError, "invalid config",
			fmt.Sprintf("config is nil"), nil)
	}
	if len(c.SigningKey) < 8 {
		return app.NewError(http.StatusInternalServerError, "invalid signing key",
			fmt.Sprintf("signing key is too short"), nil)
	}
	if c.TTLHours < 1 {
		return app.NewError(http.StatusInternalServerError, "invalid TTL hours",
			fmt.Sprintf("TTL hours are too short"), nil)
	}
	return nil
}

type tokenClaims struct {
	jwt.StandardClaims
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Nickname string    `json:"nickname"`
	Roles    []string  `json:"roles"`
}

func New(config *Config) (ports.JwtSvc, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	return &Svc{
		signingKey: config.SigningKey,
		ttl:        time.Hour * time.Duration(config.TTLHours),
	}, nil
}

func newSvc(config *Config) (*Svc, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	return &Svc{
		signingKey: config.SigningKey,
		ttl:        time.Hour * time.Duration(config.TTLHours),
	}, nil
}

func (s Svc) Generate(user d.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.ttl).Unix(),
			//IssuedAt: time.Now().Unix(),
		},
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Roles:    user.Roles.ToSlice(),
	})
	str, err := token.SignedString([]byte(s.signingKey))
	return str, err
}

func (s Svc) Parse(token string) (d.Actor, error) {

	token = strings.TrimSpace(token)

	t, err := jwt.ParseWithClaims(token,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, app.NewError(http.StatusForbidden, "invalid JWT",
					fmt.Sprintf("invalid JWT signing method parsing"), nil)
			}
			return []byte(s.signingKey), nil
		})

	if err != nil {
		return d.Actor{}, app.NewError(http.StatusForbidden, "invalid JWT",
			fmt.Sprintf("invalid JWT parsing with claims //%s//", token), err)
	}

	claims, ok := t.Claims.(*tokenClaims)
	err = validateClaims(claims)
	if !ok || err != nil {
		return d.Actor{}, app.NewError(http.StatusBadRequest, "invalid JWT, wrong claims type",
			fmt.Sprintf("token claims are not of type *tokenClaims"), nil)
	}
	// test for expired token
	if claims.ExpiresAt < time.Now().Unix() {
		return d.Actor{}, app.NewError(http.StatusUnauthorized, "token expired",
			fmt.Sprintf("token expired"), nil)
	}

	return d.NewActor(claims.ID, claims.Email, token, claims.Nickname, claims.Roles), nil
}

func validateClaims(claims *tokenClaims) error {
	if claims.ExpiresAt < time.Now().Unix() {
		return app.NewError(http.StatusUnauthorized, "token expired",
			fmt.Sprintf("token expired"), nil)
	}
	if claims.ID == uuid.Nil {
		return app.NewError(http.StatusUnauthorized, "invalid token",
			fmt.Sprintf("invalid token"), nil)
	}
	if claims.Email == "" {
		return app.NewError(http.StatusUnauthorized, "invalid token",
			fmt.Sprintf("invalid token"), nil)
	}
	if claims.Nickname == "" {
		return app.NewError(http.StatusUnauthorized, "invalid token",
			fmt.Sprintf("invalid token"), nil)
	}
	return nil
}
