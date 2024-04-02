package http

import (
	"ara-server/internal/infrastructure"
	"ara-server/internal/infrastructure/errors"
	"ara-server/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	errInvalidRequest = errors.New("invalid request")
)

type handler struct {
	infra   *infrastructure.Infrastructure
	usecase *usecase.Usecase
}

func NewHandler(infra *infrastructure.Infrastructure, usecase *usecase.Usecase) *handler {
	return &handler{
		infra:   infra,
		usecase: usecase,
	}
}

func (h *handler) RegisterHTTPHandler(router *gin.Engine) {
	router.GET("/action/last", h.initTracerContext, h.authenticate, h.HandleGetLastAction)
	router.GET("/action/available", h.initTracerContext, h.authenticate, h.HandleGetAvailableActions)
	router.GET("/action/history", h.initTracerContext, h.authenticate, h.HandleGetActionHistory)
	router.POST("/action/dispatch", h.initTracerContext, h.authenticate, h.HandleDispatchAction)

	router.GET("/actuators", h.initTracerContext, h.authenticate, h.HandleGetActuators)
	router.POST("/actuator", h.initTracerContext, h.authenticate, h.HandleInsertActuator)
	router.PATCH("/actuator", h.initTracerContext, h.authenticate, h.HandleUpdateActuator)

	router.GET("/schedules", h.initTracerContext, h.authenticate, h.HandleGetUpcomingSchedules)
	router.POST("/schedule", h.initTracerContext, h.authenticate, h.HandleCreateSchedule)
	router.POST("/scheduler/trigger", h.initTracerContext, h.authenticate, h.HandleTriggerScheduler)
	router.PATCH("/schedule", h.initTracerContext, h.authenticate, h.HandleUpdateSchedule)
	router.DELETE("/schedule", h.initTracerContext, h.authenticate, h.HandleDeleteSchedule)

	router.GET("/users", h.initTracerContext, h.authenticate, h.onlyAdmin, h.HandleGetUserInfoList)
	router.GET("/users/info", h.initTracerContext, h.authenticate, h.HandleGetUserInfo)
	router.PUT("/users", h.initTracerContext, h.authenticate, h.onlyAdmin, h.HandleUpdateUserInfo)
	router.POST("/users", h.initTracerContext, h.HandleRegisterUser)
	router.POST("/users/auth", h.initTracerContext, h.HandleLoginUser)

	router.POST("/chart", h.initTracerContext, h.authenticate, h.HandleGetSensorChart)
	router.POST("/dummy_data", h.initTracerContext, h.authenticate, h.onlyAdmin, h.HandleInsertDummyData)
}

func WriteJson(ctx *gin.Context, data interface{}, err error, httpStatusCode ...int) {
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

	if len(httpStatusCode) > 0 {
		code = httpStatusCode[0]
	}

	payload["error_code"] = code
	if code, found := errors.GetCode(err); found {
		payload["error_code"] = code
	}

	ctx.JSON(code, payload)
}
