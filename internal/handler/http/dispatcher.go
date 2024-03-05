package http

import (
	"ara-server/internal/constants"
	"ara-server/internal/usecase"

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
		ActionType: request.ActionType,
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
	actions, err := h.usecase.GetAvailableActions(c)
	if err != nil {
		WriteJson(c, nil, err)
		return
	}

	WriteJson(c, actions, nil)
}
