package http

import (
	"ara-server/internal/constants"
	"ara-server/internal/infrastructure/errors"
	"ara-server/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandleDispatchAction(c *gin.Context) error {
	var request DispatcherRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		WriteJson(c, nil, errInvalidRequest)
		return errors.Wrap(errInvalidRequest)
	}

	userID := c.GetInt(string(constants.CtxKeyUserID))

	param := usecase.DispatcherParam{
		DeviceID:   request.DeviceID,
		ActuatorID: request.ActuatorID,
		Value:      request.Value,
		ActionBy:   userID,
	}
	if err := h.usecase.DispatchAction(c, param); err != nil {
		WriteJson(c, nil, err)
		return err
	}

	WriteJson(c, nil, nil)
	return nil
}

func (h *handler) HandleGetAvailableActions(c *gin.Context) error {
	deviceIDStr := c.Query("device_id")
	if deviceIDStr == "" {
		WriteJson(c, nil, errInvalidDeviceID)
		return errInvalidDeviceID
	}

	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		WriteJson(c, nil, errInvalidDeviceID)
		return errInvalidDeviceID
	}

	actions, err := h.usecase.GetAvailableActions(c, deviceID)
	if err != nil {
		WriteJson(c, nil, err)
		return errors.Wrap(err).WithType(errors.USER)
	}

	WriteJson(c, actions, nil)
	return nil
}
