package http

import "time"

type CreateScheduleRequest struct {
	Name          string              `json:"name" binding:"required"`
	Description   string              `json:"description"`
	Actions       []DispatcherRequest `json:"actions" binding:"required"`
	ScheduledAt   time.Time           `json:"scheduled_at" binding:"required"`
	RecurringMode int                 `json:"recurring_mode"`
}
