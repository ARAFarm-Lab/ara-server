package http

import (
	"ara-server/internal/constants"
	"ara-server/internal/infrastructure/errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func (h *handler) initTracerContext(ctx *gin.Context) {
	ctx.Set(string(constants.CtxKeyCtxID), xid.New().String())
	ctx.Next()
}

func (h *handler) authenticate(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	if authorization == "" {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	segments := strings.Split(authorization, " ")
	if len(segments) != 2 {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	if strings.ToLower(segments[0]) != "bearer" {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	claims, err := h.infra.VerifyAndParseJWTToken(segments[1])
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	if userID, found := claims["user_id"]; found {
		ctx.Set(string(constants.CtxKeyUserID), userID)
	}

	ctx.Next()
}

func (h *handler) onlyAdmin(ctx *gin.Context) {
	userID := ctx.GetFloat64(string(constants.CtxKeyUserID))
	userInfo, err := h.usecase.GetUserInfo(ctx, int(userID))
	if err != nil {
		code := http.StatusInternalServerError
		if errors.IsUserError(err) {
			code = http.StatusForbidden
		}

		ctx.AbortWithError(code, err)
		return
	}

	if userInfo.Role == constants.RoleAdmin {
		ctx.Next()
		return
	}

	ctx.AbortWithStatus(http.StatusForbidden)
}
