package requests

import (
	"errors"
)

type RequestError struct {
	err error
}

func (e RequestError) message() string {
	return e.err.Error()
}

func (e RequestError) Error() string {
	return e.message()
}

type TimeoutError struct {
	RequestError
}

func NewTimeoutError(err error) TimeoutError {
	return TimeoutError{RequestError{err: err}}
}

type HTTPError struct {
	RequestError
	code int
}

type ReadBodyError struct {
	HTTPError
}

func NewHTTPError(code int, err error) HTTPError {
	return HTTPError{RequestError: RequestError{err: err}, code: code}
}

func NewReadBodyError(code int, err error) ReadBodyError {
	return ReadBodyError{NewHTTPError(code, err)}
}

func isTimeoutError(err error) bool {
	return errors.As(err, &TimeoutError{})
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	if isTimeoutError(err) {
		return true
	}
	apiErr := &HTTPError{}
	if errors.As(err, apiErr) {
		if apiErr.code >= 300 && apiErr.code < 500 {
			return false
		}
		if apiErr.code >= 500 {
			return true
		}
	}
	return false
}
