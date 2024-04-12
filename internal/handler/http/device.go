package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandleRestartDevice(ctx *gin.Context) error {
	deviceIDStr := ctx.Query("device_id")
	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		WriteJson(ctx, nil, errInvalidDeviceID)
		return errInvalidDeviceID
	}

	h.usecase.RestartDevice(ctx, deviceID)
	return nil
}
