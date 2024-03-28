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
	ID                   int                            `json:"id"`
	Name                 string                         `json:"name"`
	Description          string                         `json:"description"`
	Actions              []DispatcherParam              `json:"actions"`
	Duration             int                            `json:"duration"`
	RecurringMode        ScheduleRecurringMode          `json:"recurring_mode"`
	IsActive             bool                           `json:"is_active"`
	IsUpcomingRunCleanup bool                           `json:"is_upcoming_run_cleanup"` // indicate whether the next run is cleanup or not, only on the one-time schedule
	CleanupTime          *time.Time                     `json:"cleanup_time"`
	ScheduledAt          time.Time                      `json:"scheduled_at"`
	NextRunAt            time.Time                      `json:"next_run_at"`
	LastRunAt            time.Time                      `json:"last_run_at,omitempty"`
	LastRunStatus        constants.ActionScheduleStatus `json:"last_run_status"`
	LastError            string                         `json:"last_error"`
}

type CreateScheduleParam struct {
	Name              string
	Description       string
	Actions           []DispatcherParam
	ScheduledAt       time.Time
	DurationInMinutes int
	RecurringMode     ScheduleRecurringMode
}
