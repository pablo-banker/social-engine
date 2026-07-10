package apiErrors

const (
	// Infra / Database (100–199)
	ErrInfraBase = 100

	ErrDatabaseURLCode = ErrInfraBase + iota
	ErrConnectDBCode
	ErrPingDBCode
)
