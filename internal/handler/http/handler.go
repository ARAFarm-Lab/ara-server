package http

import (
	"ara-server/internal/infrastructure"
	"ara-server/internal/usecase"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	infra   *infrastructure.Infrastructure
	usecase *usecase.Usecase
}

var (
	errInvalidRequest = errors.New("invalid request")
)

func NewHandler(infra *infrastructure.Infrastructure, usecase *usecase.Usecase) *handler {
	return &handler{
		infra:   infra,
		usecase: usecase,
	}
}

func (h *handler) RegisterHTTPHandler(router *gin.Engine) {
	router.GET("/action/last", h.initTracerContext, h.HandleGetLastAction)
	router.GET("/action/available", h.initTracerContext, h.HandleGetAvailableActions)
	router.GET("/action/history", h.initTracerContext, h.HandleGetActionHistory)
	router.POST("/action/dispatch", h.initTracerContext, h.HandleDispatchAction)

	router.GET("/actuators", h.initTracerContext, h.HandleGetActuators)
	router.POST("/actuator", h.initTracerContext, h.HandleInsertActuator)
	router.PATCH("/actuator", h.initTracerContext, h.HandleUpdateActuator)

	router.GET("/schedules", h.initTracerContext, h.HandleGetUpcomingSchedules)
	router.POST("/schedule", h.initTracerContext, h.HandleCreateSchedule)
	router.POST("/scheduler/trigger", h.initTracerContext, h.HandleTriggerScheduler)
	router.PATCH("/schedule", h.initTracerContext, h.HandleUpdateSchedule)
	router.DELETE("/schedule", h.initTracerContext, h.HandleDeleteSchedule)

	router.GET("/users", h.authenticate, h.onlyAdmin, h.initTracerContext, h.HandleGetUserInfoList)
	router.GET("/users/info", h.authenticate, h.initTracerContext, h.HandleGetUserInfo)
	router.PUT("/users", h.authenticate, h.onlyAdmin, h.initTracerContext, h.HandleUpdateUserInfo)
	router.POST("/users", h.initTracerContext, h.HandleRegisterUser)
	router.POST("/users/auth", h.initTracerContext, h.HandleLoginUser)

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
