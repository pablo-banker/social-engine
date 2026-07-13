package models

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TokenTTL is how long an issued access token stays valid.
const TokenTTL = 7 * 24 * time.Hour

// MinSecretBytes is the minimum accepted length of JWT_SECRET. HS256 accepts
// short keys, so a weak/guessable secret would let an attacker forge tokens;
// this floor keeps the signing key infeasible to brute-force.
const MinSecretBytes = 32

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

// ValidateSecret checks that JWT_SECRET is present and long enough. Call it once
// at startup so a missing or weak secret fails fast instead of at the first
// token operation (and never reaches production silently).
func ValidateSecret() error {
	if len(os.Getenv("JWT_SECRET")) < MinSecretBytes {
		return fmt.Errorf("JWT_SECRET must be set and at least %d bytes long", MinSecretBytes)
	}
	return nil
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

// dummyPasswordHash is a valid bcrypt hash used only to burn the same CPU as a
// real password check when the account does not exist.
var dummyPasswordHash = mustBcrypt("social-engine-timing-equalizer")

func mustBcrypt(password string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return hash
}

// EqualizePasswordTiming runs a throwaway bcrypt comparison. Call it on the
// "user not found" login path so an unknown email costs the same time as a
// wrong password, closing the timing side-channel that reveals which emails
// have accounts.
func EqualizePasswordTiming(password string) {
	_ = bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(password))
}
