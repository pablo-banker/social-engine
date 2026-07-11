package apiErrors

import (
	"net/http"
)

var (
	// 1. Database & Infra
	ErrDatabaseURL = New(ErrDatabaseURLCode, "DATABASE_URL environment variable is not set", http.StatusInternalServerError)
	ErrConnectDB   = New(ErrConnectDBCode, "Error connecting to the database", http.StatusInternalServerError)
	ErrPingDB      = New(ErrPingDBCode, "Error pinging the database", http.StatusInternalServerError)
	ErrInternal    = New(ErrInternalCode, "An unexpected error occurred", http.StatusInternalServerError)

	// 2. Auth & Validação
	ErrInvalidPayload     = New(ErrInvalidPayloadCode, "The request payload is invalid", http.StatusBadRequest)
	ErrValidation         = New(ErrValidationCode, "One or more fields are invalid", http.StatusUnprocessableEntity)
	ErrInvalidCredentials = New(ErrInvalidCredentialsCode, "Invalid email or password", http.StatusUnauthorized)
	ErrEmailTaken         = New(ErrEmailTakenCode, "This email is already in use", http.StatusConflict)
	ErrUsernameTaken      = New(ErrUsernameTakenCode, "This username is already in use", http.StatusConflict)
	ErrUnauthorized       = New(ErrUnauthorizedCode, "Authentication is required", http.StatusUnauthorized)
	ErrInvalidToken       = New(ErrInvalidTokenCode, "The authentication token is invalid or expired", http.StatusUnauthorized)

	// 3. Posts
	ErrPostNotFound = New(ErrPostNotFoundCode, "Post not found", http.StatusNotFound)

	// 5. Users & Profile
	ErrUserNotFound  = New(ErrUserNotFoundCode, "User not found", http.StatusNotFound)
	ErrInvalidAvatar = New(ErrInvalidAvatarCode, "Invalid avatar selection", http.StatusUnprocessableEntity)
	ErrInvalidBanner = New(ErrInvalidBannerCode, "Invalid banner selection", http.StatusUnprocessableEntity)
)

//go:generate go run ./cmd/gen

// All is the canonical enumeration of every cataloged error, used by the
// error_code.json generator and the uniqueness regression test. New errors
// must be appended here in addition to being declared above.
var All = []*APIError{
	// 1. Database & Infra
	ErrDatabaseURL, ErrConnectDB, ErrPingDB, ErrInternal,

	// 2. Auth & Validação
	ErrInvalidPayload, ErrValidation, ErrInvalidCredentials,
	ErrEmailTaken, ErrUsernameTaken, ErrUnauthorized, ErrInvalidToken,

	// 3. Posts
	ErrPostNotFound,

	// 5. Users & Profile
	ErrUserNotFound, ErrInvalidAvatar, ErrInvalidBanner,
}
