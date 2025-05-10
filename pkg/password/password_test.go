package password

import (
	"errors"
	"music-snap/pkg/app"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:        "valid password",
			password:    "ValidPass123!",
			expectedErr: nil,
		},
		{
			name:        "too short",
			password:    "Short1!",
			expectedErr: app.NewError(http.StatusBadRequest, "invalid password", "password is too short", nil),
		},
		{
			name:        "too long",
			password:    strings.Repeat("a", 129) + "A1!",
			expectedErr: app.NewError(http.StatusBadRequest, "invalid password", "password is too long", nil),
		},
		{
			name:        "contains whitespace",
			password:    "Invalid Pass123!",
			expectedErr: app.NewError(http.StatusBadRequest, "invalid password", "password must not contain whitespace", nil),
		},
		{
			name:        "missing digit",
			password:    "InvalidPassword!",
			expectedErr: app.NewError(http.StatusBadRequest, "invalid password", "password must contain at least one digit", nil),
		},
		{
			name:        "missing special character",
			password:    "InvalidPassword123",
			expectedErr: app.NewError(http.StatusBadRequest, "invalid password", "password must contain at least one special character", nil),
		},
		{
			name:        "missing lowercase",
			password:    "INVALID123!",
			expectedErr: app.NewError(http.StatusBadRequest, "invalid password", "password must contain at least one lowercase letter", nil),
		},
		{
			name:        "missing uppercase",
			password:    "invalid123!",
			expectedErr: app.NewError(http.StatusBadRequest, "invalid password", "password must contain at least one uppercase letter", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	t.Run("empty password", func(t *testing.T) {
		_, err := HashPassword("")
		require.Error(t, err)
		assert.Equal(t, errors.New("password cannot be empty"), err)
	})

	t.Run("valid password", func(t *testing.T) {
		hash, err := HashPassword("ValidPass123!")
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.True(t, strings.HasPrefix(hash, "$argon2id$v=19$"))
	})

	t.Run("different passwords produce different hashes", func(t *testing.T) {
		hash1, err := HashPassword("Password1!")
		require.NoError(t, err)

		hash2, err := HashPassword("Password2!")
		require.NoError(t, err)

		assert.NotEqual(t, hash1, hash2)
	})
}

func TestVerifyPassword(t *testing.T) {
	validPassword := "ValidPass123!"
	invalidPassword := "InvalidPass123!"

	t.Run("valid password", func(t *testing.T) {
		hash, err := HashPassword(validPassword)
		require.NoError(t, err)

		match, err := VerifyPassword(validPassword, hash)
		require.NoError(t, err)
		assert.True(t, match)
	})

	t.Run("invalid password", func(t *testing.T) {
		hash, err := HashPassword(validPassword)
		require.NoError(t, err)

		match, err := VerifyPassword(invalidPassword, hash)
		require.NoError(t, err)
		assert.False(t, match)
	})

	t.Run("invalid hash format", func(t *testing.T) {
		_, err := VerifyPassword(validPassword, "invalid$hash$format")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid hash format parts != 6")
	})

	t.Run("unsupported algorithm", func(t *testing.T) {
		invalidHash := "$md5$v=1$m=65536,t=3,p=2$salt$hash"
		_, err := VerifyPassword(validPassword, invalidHash)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported algorithm in password hashing")
	})

	t.Run("incompatible version", func(t *testing.T) {
		invalidHash := "$argon2id$v=99$m=65536,t=3,p=2$salt$hash"
		_, err := VerifyPassword(validPassword, invalidHash)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read incompatible version in password hashing")
	})
}

func TestDecodeHash(t *testing.T) {
	validPassword := "ValidPass123!"
	validHash, err := HashPassword(validPassword)
	require.NoError(t, err)

	t.Run("valid hash", func(t *testing.T) {
		params, salt, hash, err := DecodeHash(validHash)
		require.NoError(t, err)

		assert.Equal(t, defaultParams.memory, params.memory)
		assert.Equal(t, defaultParams.iterations, params.iterations)
		assert.Equal(t, defaultParams.parallelism, params.parallelism)
		assert.Equal(t, defaultParams.saltLength, params.saltLength)
		assert.Equal(t, defaultParams.keyLength, params.keyLength)
		assert.Len(t, salt, int(defaultParams.saltLength))
		assert.Len(t, hash, int(defaultParams.keyLength))
	})

	t.Run("invalid parts count", func(t *testing.T) {
		invalidHash := "too$few$parts"
		_, _, _, err := DecodeHash(invalidHash)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid hash format parts != 6")
	})

	t.Run("invalid base64 salt", func(t *testing.T) {
		invalidHash := "$argon2id$v=19$m=65536,t=3,p=2$invalid-salt$hash"
		_, _, _, err := DecodeHash(invalidHash)
		require.Error(t, err)
	})

	t.Run("invalid base64 hash", func(t *testing.T) {
		invalidHash := "$argon2id$v=19$m=65536,t=3,p=2$salt$invalid-hash"
		_, _, _, err := DecodeHash(invalidHash)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot decode password hash")
	})
}

func TestPasswordIntegration(t *testing.T) {
	password := "SecurePass123!"

	// Hash the password
	hash, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Verify correct password
	match, err := VerifyPassword(password, hash)
	require.NoError(t, err)
	assert.True(t, match)

	// Verify incorrect password
	match, err = VerifyPassword("WrongPass123!", hash)
	require.NoError(t, err)
	assert.False(t, match)
}
