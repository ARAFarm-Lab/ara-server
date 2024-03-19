package http

import (
	"ara-server/util/log"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetActionHistory(ctx *gin.Context) {
	deviceIDStr := ctx.DefaultQuery("device_id", "-1")
	if deviceIDStr == "-1" {
		WriteJson(ctx, nil, nil, http.StatusBadRequest)
		return
	}

	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		log.Error(ctx, deviceIDStr, err, "failed to parse device ID")
		WriteJson(ctx, nil, nil, http.StatusBadRequest)
		return
	}

	result, err := h.usecase.GetActionHistories(ctx, int64(deviceID))
	if err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, result, nil)
}

func (h *handler) HandleGetLastAction(c *gin.Context) {
	deviceIDStr := c.Query("device_id")
	actuatorIDStr := c.Query("actuator_id")
	if deviceIDStr == "" || actuatorIDStr == "" {
		WriteJson(c, nil, errors.New("device_id and actuator_id are required"))
		return
	}

	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		WriteJson(c, nil, err)
		return
	}

	actuatorID, err := strconv.ParseInt(actuatorIDStr, 10, 64)
	if err != nil {
		WriteJson(c, nil, err)
		return
	}

	data, err := h.usecase.GetLastAction(deviceID, actuatorID)
	WriteJson(c, data, err)
}

func (h *handler) HandleInsertDummyData(c *gin.Context) {
	if err := h.usecase.InsertDummyData(); err != nil {
		WriteJson(c, nil, err)
		return
	}

	WriteJson(c, nil, nil)
}
