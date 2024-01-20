package http

import (
	"ara-server/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	usecase *usecase.Usecase
}

func NewHandler(usecase *usecase.Usecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

func (h *handler) RegisterHTTPHandler(router *gin.Engine) {
	router.GET("/last-action", h.initTracerContext, h.HandleGetLastAction)
	router.POST("/board/dispatch", h.initTracerContext, h.HandleDispatchAction)
	router.POST("/chart", h.initTracerContext, h.HandleGetSensorChart)
	router.POST("/scheduler/trigger", h.initTracerContext, h.HandleTriggerScheduler)
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
