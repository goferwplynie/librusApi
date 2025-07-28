package librusApi

import "errors"

type SessionExpiredError struct{}

func (e SessionExpiredError) Error() string {
	return "session expired"
}

func IsSessionExpired(e error) bool {
	var sessionExpired SessionExpiredError
	return errors.As(e, &sessionExpired)
}
