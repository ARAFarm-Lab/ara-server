package constants

type ActionType int8

const (
	ActionTypeBuiltInLED ActionType = -1
	ActionTypeRelay      ActionType = 1
)

type ActionSource uint8

const (
	ActionSourceUser       ActionSource = iota + 1 // manual trigger
	ActionSourceScheduler                          // trigger by scheduler
	ActionSourceDispatcher                         // trigger by condition defined in dispatcher
)

type ActionScheduleStatus uint8

const (
	ScheduleStatusPending ActionScheduleStatus = iota + 1
	ScheduleStatusRunning
	ScheduleStatusSuccess
	ScheduleStatusFailed
)
