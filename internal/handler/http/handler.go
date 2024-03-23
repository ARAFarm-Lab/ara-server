package http

import (
	"ara-server/internal/usecase"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	usecase *usecase.Usecase
}

var (
	errInvalidRequestBody = errors.New("invalid request body")
)

func NewHandler(usecase *usecase.Usecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

func (h *handler) RegisterHTTPHandler(router *gin.Engine) {
	router.GET("/action/last", h.initTracerContext, h.HandleGetLastAction)
	router.GET("/action/available", h.initTracerContext, h.HandleGetAvailableActions)
	router.GET("/action/history", h.initTracerContext, h.HandleGetActionHistory)
	router.POST("/action/dispatch", h.initTracerContext, h.HandleDispatchAction)

	router.GET("/actuators", h.initTracerContext, h.HandleGetActiveActuators)
	router.POST("/actuator", h.initTracerContext, h.HandleInsertActuator)
	router.PATCH("/actuator", h.initTracerContext, h.HandleUpdateActuator)

	router.GET("/schedules", h.initTracerContext, h.HandleGetUpcomingSchedules)
	router.POST("/schedule", h.initTracerContext, h.HandleCreateSchedule)
	router.POST("/scheduler/trigger", h.initTracerContext, h.HandleTriggerScheduler)
	router.PATCH("/schedule", h.initTracerContext, h.HandleUpdateSchedule)
	router.DELETE("/schedule", h.initTracerContext, h.HandleDeleteSchedule)

	router.POST("/chart", h.initTracerContext, h.HandleGetSensorChart)
	router.POST("/dummy_data", h.initTracerContext, h.HandleInsertDummyData)
}

func WriteJson(ctx *gin.Context, data interface{}, err error, statusCode ...int) {
	payload := map[string]interface{}{
		"is_success": true,
	}
	code := http.StatusOK
	if data != nil {
		payload["data"] = data
	}

	if err != nil {
		code = http.StatusInternalServerError
		payload["is_success"] = false
		payload["error"] = err.Error()
	}

	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	payload["code"] = code

	ctx.JSON(code, payload)
}
