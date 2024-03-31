package http

import (
	"ara-server/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetUpcomingSchedules(ctx *gin.Context) {
	schedules, err := h.usecase.GetUpcomingSchedules(ctx)
	if err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, schedules, nil)
}

func (h *handler) HandleCreateSchedule(ctx *gin.Context) {
	var request CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		WriteJson(ctx, nil, errInvalidRequest, http.StatusBadRequest)
		return
	}

	actions := make([]usecase.DispatcherParam, len(request.Actions))
	for i, action := range request.Actions {
		actions[i] = usecase.DispatcherParam{
			DeviceID:   action.DeviceID,
			ActuatorID: action.ActuatorID,
			Value:      action.Value,
		}
	}

	param := usecase.CreateScheduleParam{
		Name:              request.Name,
		Description:       request.Description,
		Actions:           actions,
		ScheduledAt:       request.ScheduledAt,
		DurationInMinutes: request.Duration,
		RecurringMode:     usecase.ScheduleRecurringMode(request.RecurringMode),
	}
	if err := h.usecase.CreateSchedule(ctx, param); err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, nil, nil, http.StatusCreated)
}

func (h *handler) HandleDeleteSchedule(ctx *gin.Context) {
	scheduleIDStr := ctx.DefaultQuery("id", "-1")
	if scheduleIDStr == "-1" {
		WriteJson(ctx, nil, nil, http.StatusBadRequest)
		return
	}

	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		WriteJson(ctx, nil, nil, http.StatusBadRequest)
		return
	}

	if err := h.usecase.DeleteSchedule(ctx, scheduleID); err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, nil, nil, http.StatusAccepted)
}

func (h *handler) HandleUpdateSchedule(ctx *gin.Context) {
	var request usecase.ActionSchedule
	if err := ctx.ShouldBindJSON(&request); err != nil {
		WriteJson(ctx, nil, errInvalidRequest, http.StatusBadRequest)
		return
	}

	if err := h.usecase.UpdateSchedule(ctx, request); err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, nil, nil, http.StatusCreated)
}

func (h *handler) HandleTriggerScheduler(ctx *gin.Context) {
	go func() {
		h.usecase.DispatchScheduler(ctx)
	}()

	WriteJson(ctx, nil, nil)
}
