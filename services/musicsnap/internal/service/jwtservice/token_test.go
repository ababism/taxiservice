package jwtservice

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	d "music-snap/services/musicsnap/internal/domain"
	"testing"
	"time"
)

const singKey = "30f84aa3dd8fdc9634e6e4697b2d3ce3"

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectedErr string
	}{
		{
			name:        "nil config",
			config:      nil,
			expectedErr: "config is nil",
		},
		{
			name: "short signing key",
			config: &Config{
				SigningKey: "short",
				TTLHours:   1,
			},
			expectedErr: "signing key is too short",
		},
		{
			name: "invalid TTL",
			config: &Config{
				SigningKey: "valid-signing-key",
				TTLHours:   0,
			},
			expectedErr: "TTL hours are too short",
		},
		{
			name: "valid config",
			config: &Config{
				SigningKey: singKey,
				TTLHours:   1,
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectedErr string
	}{
		{
			name:        "nil config",
			config:      nil,
			expectedErr: "config is nil",
		},
		{
			name: "valid config",
			config: &Config{
				SigningKey: singKey,
				TTLHours:   1,
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := newSvc(tt.config)
			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, svc)
			} else {
				require.NoError(t, err)
				require.NotNil(t, svc)
				assert.Equal(t, tt.config.SigningKey, svc.signingKey)
				assert.Equal(t, time.Hour*time.Duration(tt.config.TTLHours), svc.ttl)
			}
		})
	}
}

func TestSvc_Generate(t *testing.T) {
	validConfig := &Config{
		SigningKey: singKey,
		TTLHours:   1,
	}

	svc, err := newSvc(validConfig)
	require.NoError(t, err)
	require.NotNil(t, svc)
	//profile := domain.Profile{
	//	ID:            uuid.New(),
	//	Nickname:      "testuser",
	//	AvatarURL:     "http://example.com/avatar.png",
	//	BackgroundURL: "http://example.com/bg.png",
	//	Bio:           "Test bio",
	//}
	//testUser := domain.User{
	//	Profile:      profile,
	//	Email:        "test@example.com",
	//	PasswordHash: "hashedpassword",
	//	Roles:        domain.NewRoles([]string{domain.UserRole, domain.RegisteredRole}),
	//}
	testProfile := d.Profile{
		ID:            uuid.New(),
		Nickname:      "testuser",
		AvatarURL:     "http://example.com/avatar.png",
		BackgroundURL: "http://example.com/bg.png",
		Bio:           "Test bio",
	}
	testUser := d.User{
		Profile:      testProfile,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Roles:        d.NewRoles([]string{d.UserRole, d.RegisteredRole}),
	}

	t.Run("successful token generation", func(t *testing.T) {
		token, err := svc.Generate(testUser)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify the token can be parsed back
		parsedToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(validConfig.SigningKey), nil
		})
		require.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(*tokenClaims)
		require.True(t, ok)
		assert.Equal(t, testUser.ID, claims.ID)
		assert.Equal(t, testUser.Email, claims.Email)
		assert.Equal(t, testUser.Nickname, claims.Nickname)
		assert.ElementsMatch(t, testUser.Roles.ToSlice(), claims.Roles)
		assert.WithinDuration(t, time.Now().Add(svc.ttl), time.Unix(claims.ExpiresAt, 0), time.Second)
	})
}

func TestSvc_Parse(t *testing.T) {
	validConfig := &Config{
		SigningKey: singKey,
		TTLHours:   1,
	}

	svc, err := New(validConfig)
	require.NoError(t, err)
	require.NotNil(t, svc)

	testProfile := d.Profile{
		ID:            uuid.New(),
		Nickname:      "testuser",
		AvatarURL:     "http://example.com/avatar.png",
		BackgroundURL: "http://example.com/bg.png",
		Bio:           "Test bio",
	}
	testUser := d.User{
		Profile:      testProfile,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Roles:        d.NewRoles([]string{d.UserRole, d.RegisteredRole}),
	}

	t.Run("successful token parsing", func(t *testing.T) {
		token, err := svc.Generate(testUser)
		require.NoError(t, err)

		actor, err := svc.Parse(token)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, actor.ID)
		assert.Equal(t, testUser.Email, actor.Mail)
		assert.Equal(t, testUser.Nickname, actor.Nickname)
		assert.ElementsMatch(t, testUser.Roles.ToSlice(), actor.GetRoles())
		assert.Equal(t, token, actor.Jwt)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := svc.Parse("invalid.token.string")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid JWT")
	})

	t.Run("wrong signing key", func(t *testing.T) {
		// Create a token with different signing key
		otherConfig := &Config{
			SigningKey: "different-signing-key",
			TTLHours:   1,
		}
		otherSvc, err := New(otherConfig)
		require.NoError(t, err)

		token, err := otherSvc.Generate(testUser)
		require.NoError(t, err)

		_, err = svc.Parse(token)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid JWT")
	})

	t.Run("invalid claims type", func(t *testing.T) {
		// Create a token with different claims type
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, err := token.SignedString([]byte(validConfig.SigningKey))
		require.NoError(t, err)

		_, err = svc.Parse(tokenString)
		fmt.Printf("error: %v\n", err)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "wrong claims type")
	})

	t.Run("static token parse", func(t *testing.T) {
		temp, tempErr := svc.Parse(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAyNzE2NzgsImlkIjoiN2Q5YjJjZWMtZGFjZS00MGFlLTliNjctZDViMDMxNmNlYThmIiwiZW1haWwiOiJraXNlbGV2LmNvbnRhY3RAZ21haWwuY29tIiwibmlja25hbWUiOiJhYmFiaXNtIiwicm9sZXMiOlsiYWRtaW4iLCJ1c2VyIl19.ojniKjrxbUb0VOH-YiZBWBVCoD5NHbdbHZxIRyWzJfU",
		)

		assert.NoError(t, tempErr)
		assert.Equal(t, "7d9b2cec-dace-40ae-9b67-d5b0316cea8f", temp.ID.String())
	})
}
