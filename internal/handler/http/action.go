package http

import (
	"ara-server/internal/infrastructure/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetActionHistory(ctx *gin.Context) error {
	deviceIDStr := ctx.DefaultQuery("device_id", "")
	if deviceIDStr == "" {
		err := errors.New("invalid device id").WithType(errors.USER)
		WriteJson(ctx, nil, err, http.StatusBadRequest)
		return err
	}

	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		WriteJson(ctx, nil, err, http.StatusBadRequest)
		return errors.Wrap(err).WithType(errors.USER)
	}

	result, err := h.usecase.GetActionHistories(ctx, int64(deviceID))
	if err != nil {
		WriteJson(ctx, nil, err)
		return err
	}

	WriteJson(ctx, result, nil)
	return nil
}

func (h *handler) HandleGetLastAction(c *gin.Context) error {
	deviceIDStr := c.Query("device_id")
	actuatorIDStr := c.Query("actuator_id")
	if deviceIDStr == "" || actuatorIDStr == "" {
		err := errors.New("device_id and actuator_id are required").WithType(errors.USER)
		WriteJson(c, nil, err)
		return err
	}

	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		WriteJson(c, nil, err)
		return errors.Wrap(err).WithType(errors.USER)
	}

	actuatorID, err := strconv.ParseInt(actuatorIDStr, 10, 64)
	if err != nil {
		WriteJson(c, nil, err)
		errors.Wrap(err).WithType(errors.USER)
	}

	data, err := h.usecase.GetLastAction(deviceID, actuatorID)
	WriteJson(c, data, err)
	return nil
}

func (h *handler) HandleInsertDummyData(c *gin.Context) error {
	if err := h.usecase.InsertDummyData(); err != nil {
		WriteJson(c, nil, err)
		return err
	}

	WriteJson(c, nil, nil)
	return nil
}
