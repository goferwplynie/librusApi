package librusApi

import "errors"

// error thrown when session expired. librus api doesnt explicitely say that but when client was logged in and gets back Access Denied it assumes that session expired.
type SessionExpiredError struct{}

func (e SessionExpiredError) Error() string {
	return "session expired"
}

func IsSessionExpired(e error) bool {
	var sessionExpired SessionExpiredError
	return errors.As(e, &sessionExpired)
}
