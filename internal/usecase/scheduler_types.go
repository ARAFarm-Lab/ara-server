package usecase

import (
	"ara-server/internal/constants"
	"time"
)

type ScheduleRecurringMode uint8

const (
	RecurringModeNone     ScheduleRecurringMode = 0
	RecurringModeMinutely ScheduleRecurringMode = 1
	RecurringModeHourly   ScheduleRecurringMode = 2
	RecurringModeDaily    ScheduleRecurringMode = 3
)

type ActionSchedule struct {
	ID            int                            `json:"id"`
	Name          string                         `json:"name"`
	Description   string                         `json:"description"`
	Actions       []DispatcherParam              `json:"actions"`
	RecurringMode ScheduleRecurringMode          `json:"recurring_mode"`
	IsActive      bool                           `json:"is_active"`
	ScheduledAt   time.Time                      `json:"scheduled_at"`
	LastRunAt     time.Time                      `json:"last_run_at"`
	LastRunStatus constants.ActionScheduleStatus `json:"last_run_status"`
	LastError     string                         `json:"last_error"`
}

type CreateScheduleParam struct {
	Name          string
	Description   string
	Actions       []DispatcherParam
	ScheduledAt   time.Time
	RecurringMode ScheduleRecurringMode
}
