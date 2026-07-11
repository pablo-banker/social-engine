package apiErrors

const (
	// Infra / Database (100–199)
	ErrInfraBase = 100

	ErrDatabaseURLCode = ErrInfraBase + iota
	ErrConnectDBCode
	ErrPingDBCode
	ErrInternalCode
)

const (
	// Auth & Validação (200–299)
	ErrAuthBase = 200

	ErrInvalidPayloadCode = ErrAuthBase + iota
	ErrValidationCode
	ErrInvalidCredentialsCode
	ErrEmailTakenCode
	ErrUsernameTakenCode
	ErrUnauthorizedCode
	ErrInvalidTokenCode
)

const (
	// Posts (300–399)
	ErrPostsBase = 300

	ErrPostNotFoundCode = ErrPostsBase + iota
)

const (
	// Users & Profile (500–599)
	ErrUsersBase = 500

	ErrUserNotFoundCode = ErrUsersBase + iota
	ErrInvalidAvatarCode
	ErrInvalidBannerCode
)
