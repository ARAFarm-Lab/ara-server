package errors

import (
	"ara-server/internal/constants"
	"errors"
)

type errorType uint8

const (
	USER   errorType = 1
	SYSTEM errorType = 2
)

type customError struct {
	err       error
	errorType errorType
	code      constants.ErrorCode
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

func GetCode(err error) (constants.ErrorCode, bool) {
	customError, ok := err.(customError)
	if !ok {
		return 0, false
	}

	return customError.code, customError.code != 0
}

func (e customError) WithType(errorType errorType) customError {
	e.errorType = errorType
	return e
}

func (e customError) WithCode(code constants.ErrorCode) customError {
	e.code = code
	return e
}

func Wrap(err error) customError {
	return customError{
		err:       err,
		errorType: SYSTEM,
	}
}

func IsUserError(err error) bool {
	customError, ok := err.(customError)
	if !ok {
		return false
	}

	return customError.errorType == USER
}
