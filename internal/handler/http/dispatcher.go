package http

import (
	"ara-server/internal/constants"
	"ara-server/internal/usecase"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandleDispatchAction(c *gin.Context) {
	var request DispatcherRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		WriteJson(c, nil, err)
		return
	}

	param := usecase.DispatcherParam{
		DeviceID:   request.DeviceID,
		ActuatorID: request.ActuatorID,
		Value:      request.Value,
		ActionBy:   constants.ActionSourceUser,
	}
	if err := h.usecase.DispatchAction(c, param); err != nil {
		WriteJson(c, nil, err)
		return
	}

	WriteJson(c, nil, nil)
}

func (h *handler) HandleGetAvailableActions(c *gin.Context) {
	deviceIDStr := c.Query("device_id")
	if deviceIDStr == "" {
		WriteJson(c, nil, errors.New("device_id is missing"))
		return
	}

	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		WriteJson(c, nil, err)
		return
	}

	actions, err := h.usecase.GetAvailableActions(c, deviceID)
	if err != nil {
		WriteJson(c, nil, err)
		return
	}

	WriteJson(c, actions, nil)
}
