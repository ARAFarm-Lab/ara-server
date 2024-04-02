package constants

import "errors"

var (
	// schedule error
	ErrorScheduleNotFound = errors.New("schedule not found")

	// actuator error
	ErrorActuatorNotFound = errors.New("actuator not found")
)

// custom error codes for client
type ErrorCode int

const (
	ErrorCodeUserNotFound ErrorCode = 103
	ErrorCodeUserExists   ErrorCode = 104
)
