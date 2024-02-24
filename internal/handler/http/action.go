package http

import (
	"ara-server/internal/constants"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
