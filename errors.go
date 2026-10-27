package librus

import (
	"errors"
	"fmt"
)

var (
	// ErrNotLoggedIn is returned when a resource is requested before a successful Login.
	ErrNotLoggedIn = errors.New("librus: not logged in, call Login first")
	// ErrLoginFailed is returned by Login when Librus rejects the credentials.
	ErrLoginFailed = errors.New("librus: login failed, check login and password")
)

// APIError is returned when Librus responds with status code >= 400.
type APIError struct {
	StatusCode int
	// Code is Librus's machine readable error code (e.g. "LuckyNumberIsNotActive"). Often empty.
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

func (e *APIError) Error() string {
	msg := fmt.Sprintf("librus: api returned %d", e.StatusCode)
	if e.Code != "" {
		msg += " " + e.Code
	}
	if e.Message != "" {
		msg += ": " + e.Message
	}
	return msg
}

// sessionExpired reports whether the error means the session is gone. Librus
// answers 401 for other reasons too, e.g. "Insufficient scopes" on Messages.
func (e *APIError) sessionExpired() bool {
	return e.StatusCode == 401 && e.Message == "Request is denied"
}
