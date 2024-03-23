package constants

import "errors"

var (
	// schedule error
	ErrorScheduleNotFound = errors.New("schedule not found")

	// actuator error
	ErrorActuatorNotFound = errors.New("actuator not found")
)
