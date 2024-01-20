package http

import "github.com/gin-gonic/gin"

func (h *handler) HandleTriggerScheduler(ctx *gin.Context) {
	go func() {
		h.usecase.DispatchScheduler(ctx)
	}()

	WriteJson(ctx, nil, nil)
}
