package apiErrors

import (
	"net/http"
)

var (
	// 1. Database & Infra
	ErrDatabaseURL = New(ErrDatabaseURLCode, "DATABASE_URL environment variable is not set", http.StatusInternalServerError)
	ErrConnectDB   = New(ErrConnectDBCode, "Error connecting to the database", http.StatusInternalServerError)
	ErrPingDB      = New(ErrPingDBCode, "Error pinging the database", http.StatusInternalServerError)
)

//go:generate go run ./cmd/gen

// All is the canonical enumeration of every cataloged error, used by the
// error_code.json generator and the uniqueness regression test. New errors
// must be appended here in addition to being declared above.
var All = []*APIError{
	// 1. Database & Infra
	ErrDatabaseURL, ErrConnectDB, ErrPingDB,
}
