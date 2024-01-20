package http

import (
	"ara-server/internal/constants"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func (h *handler) initTracerContext(ctx *gin.Context) {
	ctx.Set(string(constants.CtxKeyCtxID), xid.New().String())
	ctx.Next()
}
