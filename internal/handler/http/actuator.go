package http

import (
	"ara-server/internal/constants"
	"ara-server/internal/usecase"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetActuators(ctx *gin.Context) {
	deviceIDStr := ctx.Query("device_id")
	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	result, err := h.usecase.GetActuators(ctx, deviceID)
	if err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, result, nil)
}

func (h *handler) HandleInsertActuator(ctx *gin.Context) {
	var actuator usecase.Actuator
	if err := ctx.ShouldBindJSON(&actuator); err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	if err := h.usecase.InsertActuator(ctx, actuator); err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, nil, nil)
}

func (h *handler) HandleUpdateActuator(ctx *gin.Context) {
	var actuator usecase.Actuator
	if err := ctx.ShouldBindJSON(&actuator); err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	ctxHandler := context.Background()
	for key, value := range ctx.Copy().Keys {
		ctxHandler = context.WithValue(ctxHandler, constants.ContextKey(key), value)
	}

	if err := h.usecase.UpdateActuator(ctx, actuator); err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, nil, nil)
}
