package http

import (
	"ara-server/internal/constants"
	"ara-server/internal/infrastructure"
	"ara-server/internal/infrastructure/errors"
	"ara-server/internal/infrastructure/metric"
	"ara-server/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rs/xid"
)

var (
	errInvalidRequest  = errors.New("invalid request").WithType(errors.USER)
	errInvalidDeviceID = errors.New("invalid device_id").WithType(errors.USER)
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
	router.GET("/action/last", h.authenticate, h.handlerWrapper(h.HandleGetLastAction))
	router.GET("/action/available", h.authenticate, h.handlerWrapper(h.HandleGetAvailableActions))
	router.GET("/action/history", h.authenticate, h.handlerWrapper(h.HandleGetActionHistory))
	router.POST("/action/dispatch", h.authenticate, h.handlerWrapper(h.HandleDispatchAction))

	router.GET("/actuators", h.authenticate, h.handlerWrapper(h.HandleGetActuators))
	router.POST("/actuator", h.authenticate, h.handlerWrapper(h.HandleInsertActuator))
	router.PATCH("/actuator", h.authenticate, h.handlerWrapper(h.HandleUpdateActuator))

	router.GET("/schedules", h.authenticate, h.handlerWrapper(h.HandleGetUpcomingSchedules))
	router.POST("/schedule", h.authenticate, h.handlerWrapper(h.HandleCreateSchedule))
	router.POST("/scheduler/trigger", h.authenticate, h.handlerWrapper(h.HandleTriggerScheduler))
	router.PATCH("/schedule", h.authenticate, h.handlerWrapper(h.HandleUpdateSchedule))
	router.DELETE("/schedule", h.authenticate, h.handlerWrapper(h.HandleDeleteSchedule))

	router.GET("/users", h.authenticate, h.onlyAdmin, h.handlerWrapper(h.HandleGetUserInfoList))
	router.GET("/users/info", h.authenticate, h.handlerWrapper(h.HandleGetUserInfo))
	router.PUT("/users", h.authenticate, h.onlyAdmin, h.handlerWrapper(h.HandleUpdateUserInfo))
	router.POST("/users", h.handlerWrapper(h.HandleRegisterUser))
	router.POST("/users/auth", h.handlerWrapper(h.HandleLoginUser))

	router.POST("/chart", h.authenticate, h.handlerWrapper(h.HandleGetSensorChart))
	router.POST("/dummy_data", h.authenticate, h.onlyAdmin, h.handlerWrapper(h.HandleInsertDummyData))

	router.GET("/_ara-iot/metric", gin.WrapH(promhttp.Handler()))
}

func (h *handler) handlerWrapper(handler func(ctx *gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Initialize tracer context
		ctx.Set(string(constants.CtxKeyCtxID), xid.New().String())
		ctx.Next()

		err := handler(ctx)

		isSuccess := "true"
		if err != nil {
			isSuccess = "false"
		}
		h.infra.Metric.PushCounter(metric.HTTPResponse, map[string]string{
			"method":     ctx.Request.Method,
			"path":       ctx.Request.URL.Path,
			"is_success": isSuccess,
		})
	}
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
