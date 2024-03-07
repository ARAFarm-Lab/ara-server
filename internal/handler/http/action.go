package http

import (
	"ara-server/internal/constants"
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
	actionTypeStr := c.Query("action_type")
	if deviceIDStr == "" || actionTypeStr == "" {
		WriteJson(c, nil, errors.New("device_id and sensor_type are required"))
		return
	}

	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		WriteJson(c, nil, err)
		return
	}

	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil {
		WriteJson(c, nil, err)
		return
	}

	data, err := h.usecase.GetLastAction(deviceID, constants.ActionType(actionType))
	WriteJson(c, data, err)
}

func (h *handler) HandleInsertDummyData(c *gin.Context) {
	if err := h.usecase.InsertDummyData(); err != nil {
		WriteJson(c, nil, err)
		return
	}

	WriteJson(c, nil, nil)
}
