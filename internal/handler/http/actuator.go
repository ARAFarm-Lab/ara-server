package http

import (
	"ara-server/internal/usecase"

	"github.com/gin-gonic/gin"
)

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
