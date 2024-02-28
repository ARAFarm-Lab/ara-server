package http

import (
	"ara-server/internal/usecase"
	"ara-server/util/log"
	"net/http"

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
		log.Error(ctx, nil, err, "failed to bind json")
		WriteJson(ctx, nil, errInvalidRequestBody, http.StatusBadRequest)
		return
	}

	actions := make([]usecase.DispatcherParam, len(request.Actions))
	for i, action := range request.Actions {
		actions[i] = usecase.DispatcherParam(action)
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

func (h *handler) HandleTriggerScheduler(ctx *gin.Context) {
	go func() {
		h.usecase.DispatchScheduler(ctx)
	}()

	WriteJson(ctx, nil, nil)
}
