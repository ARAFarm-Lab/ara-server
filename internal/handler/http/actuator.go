package http

import (
	"ara-server/internal/constants"
	"ara-server/internal/infrastructure/errors"
	"ara-server/internal/usecase"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetActuators(ctx *gin.Context) error {
	deviceIDStr := ctx.Query("device_id")
	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		WriteJson(ctx, nil, errInvalidDeviceID)
		return errInvalidDeviceID
	}

	result, err := h.usecase.GetActuators(ctx, deviceID)
	if err != nil {
		WriteJson(ctx, nil, err)
		return err
	}

	WriteJson(ctx, result, nil)
	return nil
}

func (h *handler) HandleInsertActuator(ctx *gin.Context) error {
	var actuator usecase.Actuator
	if err := ctx.ShouldBindJSON(&actuator); err != nil {
		WriteJson(ctx, nil, err)
		return errors.Wrap(err).WithType(errors.USER)
	}

	if err := h.usecase.InsertActuator(ctx, actuator); err != nil {
		WriteJson(ctx, nil, err)
		return err
	}

	WriteJson(ctx, nil, nil)
	return nil
}

func (h *handler) HandleUpdateActuator(ctx *gin.Context) error {
	var actuator usecase.Actuator
	if err := ctx.ShouldBindJSON(&actuator); err != nil {
		WriteJson(ctx, nil, err)
		return errors.Wrap(err).WithType(errors.USER)
	}

	ctxHandler := context.Background()
	for key, value := range ctx.Copy().Keys {
		ctxHandler = context.WithValue(ctxHandler, constants.ContextKey(key), value)
	}

	if err := h.usecase.UpdateActuator(ctx, actuator); err != nil {
		WriteJson(ctx, nil, err)
		return err
	}

	WriteJson(ctx, nil, nil)
	return nil
}
