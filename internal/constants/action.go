package constants

type ActionType int8

const (
	ActionTypeBuiltInLED ActionType = -1
	ActionTypeRelay      ActionType = 1
)

const (
	ActionSourceScheduler  = -1
	ActionSourceDispatcher = -2
)

type ActionScheduleStatus uint8

const (
	ScheduleStatusPending ActionScheduleStatus = iota + 1
	ScheduleStatusRunning
	ScheduleStatusSuccess
	ScheduleStatusFailed
)
