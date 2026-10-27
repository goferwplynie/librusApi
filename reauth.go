package librus

import "fmt"

// WithReauth calls method and, if it fails with SessionExpiredError, calls loginFunc
// and retries up to retries times.
func WithReauth[T any](
	method func() (T, error),
	loginFunc func() error,
	retries int,
) (T, error) {
	var zeroVal T
	for range retries + 1 {
		resp, err := method()
		if err == nil {
			return resp, nil
		}

		if IsSessionExpired(err) {
			err := loginFunc()
			if err != nil {
				return resp, fmt.Errorf("error in login func: %v", err)
			}
			continue
		}
		return zeroVal, err
	}
	return zeroVal, fmt.Errorf("reauth failed")
}
