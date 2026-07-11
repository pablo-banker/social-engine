package models

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TokenTTL is how long an issued access token stays valid.
const TokenTTL = 7 * 24 * time.Hour

// ErrMissingSecret is returned when the JWT signing secret is not configured.
var ErrMissingSecret = errors.New("JWT_SECRET environment variable is not set")

// ErrInvalidToken is returned when a token fails signature or claims validation.
var ErrInvalidToken = errors.New("authentication token is invalid or expired")

func secret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return nil, ErrMissingSecret
	}
	return []byte(s), nil
}

// GenerateToken issues a signed HS256 token whose subject is the user ID.
func GenerateToken(userID uuid.UUID) (string, error) {
	key, err := secret()
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(TokenTTL)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
}

// ParseToken validates a token and returns the user ID stored in its subject.
func ParseToken(tokenString string) (uuid.UUID, error) {
	key, err := secret()
	if err != nil {
		return uuid.Nil, err
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return key, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	return userID, nil
}

// HashPassword returns the bcrypt hash of the given plaintext password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword reports whether the plaintext password matches the bcrypt hash.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
