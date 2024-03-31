package errors

import "errors"

type errorType uint8

const (
	USER   errorType = 1
	SYSTEM errorType = 2
)

type customError struct {
	err       error
	errorType errorType
}

func (e customError) Error() string {
	return e.err.Error()
}

func New(message string) customError {
	return customError{
		err:       errors.New(message),
		errorType: SYSTEM,
	}
}

func Wrap(err error) customError {
	return customError{
		err:       err,
		errorType: SYSTEM,
	}
}

func (e customError) WithType(errorType errorType) customError {
	e.errorType = errorType
	return e
}

func IsUserError(err error) bool {
	customError, ok := err.(customError)
	if !ok {
		return false
	}

	return customError.errorType == USER
}
