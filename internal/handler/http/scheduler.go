package http

import (
	"ara-server/internal/infrastructure/errors"
	"ara-server/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetUpcomingSchedules(ctx *gin.Context) error {
	schedules, err := h.usecase.GetUpcomingSchedules(ctx)
	if err != nil {
		WriteJson(ctx, nil, err)
		return err
	}

	WriteJson(ctx, schedules, nil)
	return nil
}

func (h *handler) HandleCreateSchedule(ctx *gin.Context) error {
	var request CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		WriteJson(ctx, nil, errInvalidRequest, http.StatusBadRequest)
		return errInvalidRequest
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
		return err
	}

	WriteJson(ctx, nil, nil, http.StatusCreated)
	return nil
}

func (h *handler) HandleDeleteSchedule(ctx *gin.Context) error {
	scheduleIDStr := ctx.DefaultQuery("id", "-1")
	if scheduleIDStr == "-1" {
		WriteJson(ctx, nil, errInvalidDeviceID, http.StatusBadRequest)
		return errInvalidDeviceID
	}

	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		err := errors.New("invalid schedule_id").WithType(errors.USER)
		WriteJson(ctx, nil, err, http.StatusBadRequest)
		return err
	}

	if err := h.usecase.DeleteSchedule(ctx, scheduleID); err != nil {
		WriteJson(ctx, nil, err)
		return err
	}

	WriteJson(ctx, nil, nil, http.StatusAccepted)
	return nil
}

func (h *handler) HandleUpdateSchedule(ctx *gin.Context) error {
	var request usecase.ActionSchedule
	if err := ctx.ShouldBindJSON(&request); err != nil {
		WriteJson(ctx, nil, errInvalidRequest, http.StatusBadRequest)
		return errInvalidRequest
	}

	if err := h.usecase.UpdateSchedule(ctx, request); err != nil {
		WriteJson(ctx, nil, err)
		return err
	}

	WriteJson(ctx, nil, nil, http.StatusCreated)
	return nil
}

func (h *handler) HandleTriggerScheduler(ctx *gin.Context) error {
	go func() {
		h.usecase.DispatchScheduler(ctx)
	}()

	WriteJson(ctx, nil, nil)
	return nil
}
