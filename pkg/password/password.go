package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"music-snap/pkg/app"
	"net/http"
	"regexp"
	"strings"
)

// Parameters for Argon2 hashing
type argonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// Default parameters (you can adjust these based on your security requirements)
var defaultParams = &argonParams{
	memory:      64 * 1024, // 64 MB
	iterations:  3,
	parallelism: 2,
	saltLength:  16, // 128-bit salt
	keyLength:   32, // 256-bit key
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return app.NewError(http.StatusBadRequest, "invalid password",
			fmt.Sprintf("password is too short"), nil)
	}
	if len(password) > 128 {
		return app.NewError(http.StatusBadRequest, "invalid password",
			fmt.Sprintf("password is too long"), nil)
	}

	// whitespace
	if strings.ContainsAny(password, " \t") {
		return app.NewError(http.StatusBadRequest, "invalid password",
			fmt.Sprintf("password must not contain whitespace"), nil)
	}
	if strings.ContainsAny(password, " \n") {
		return app.NewError(http.StatusBadRequest, "invalid password",
			fmt.Sprintf("password must not contain whitespace"), nil)
	}

	// block of code to check that password contains at least one digit letter uppercase letter and special character
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return app.NewError(http.StatusBadRequest, "invalid password",
			fmt.Sprintf("password must contain at least one digit"), nil)
	}
	if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		return app.NewError(http.StatusBadRequest, "invalid password",
			fmt.Sprintf("password must contain at least one special character"), nil)
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return app.NewError(http.StatusBadRequest, "invalid password",
			fmt.Sprintf("password must contain at least one lowercase letter"), nil)
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return app.NewError(http.StatusBadRequest, "invalid password",
			fmt.Sprintf("password must contain at least one uppercase letter"), nil)
	}

	return nil
}

// HashPassword generates a secure hash of a password using Argon2
func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password cannot be empty")

	}

	// Generate a cryptographically secure random salt
	salt := make([]byte, defaultParams.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Generate the hash using Argon2id (variant that's resistant to both GPU and side-channel attacks)
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		defaultParams.iterations,
		defaultParams.memory,
		defaultParams.parallelism,
		defaultParams.keyLength,
	)

	// Encode the hash and parameters into a string for storage
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		defaultParams.memory,
		defaultParams.iterations,
		defaultParams.parallelism,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

// VerifyPassword compare plain-text password with hashed password
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Extract the parameters, salt, and derived key from the encoded hash
	params, salt, hash, err := DecodeHash(encodedHash)
	if err != nil {
		return false, app.NewError(http.StatusInternalServerError, "checking password", "checking password", err)
	}

	// Derive the key from the password using the same parameters
	derivedKey := argon2.IDKey(
		[]byte(password),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength,
	)

	// Constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare(hash, derivedKey) == 1 {
		return true, nil
	}

	return false, nil
}

// DecodeHash decodes the stored hash into its components
func DecodeHash(encodedHash string) (*argonParams, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, app.NewError(http.StatusInternalServerError, "password hashing error", "invalid hash format parts != 6", nil)
	}

	// Check the algorithm
	if parts[1] != "argon2id" {
		return nil, nil, nil, app.NewError(http.StatusInternalServerError, "password hashing error", "unsupported algorithm in password hashing", nil)
	}

	// Parse version
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, app.NewError(http.StatusInternalServerError, "password hashing error", "cannot read incompatible version in password hashing", nil)
	}

	// Parse parameters
	params := &argonParams{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism); err != nil {
		return nil, nil, nil, app.NewError(http.StatusInternalServerError, "password hashing error", "cannot read argon parameters in password hashing", err)
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.saltLength = uint32(len(salt))

	// Decode hash
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, app.NewError(http.StatusInternalServerError, "password hashing error", "cannot decode password hash", err)
	}
	params.keyLength = uint32(len(hash))

	return params, salt, hash, nil
}
