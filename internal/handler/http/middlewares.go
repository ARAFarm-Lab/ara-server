package http

import (
	"ara-server/internal/constants"
	"ara-server/internal/infrastructure/errors"
	"ara-server/util/log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var forbiddenPayload = map[string]interface{}{
	"message":    "Access Denied",
	"error_code": http.StatusForbidden,
}

func (h *handler) initTracerContext(ctx *gin.Context) {
	ctx.Set(string(constants.CtxKeyCtxID), xid.New().String())
	ctx.Next()
}

func (h *handler) authenticate(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	if authorization == "" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, forbiddenPayload)
		return
	}

	segments := strings.Split(authorization, " ")
	if len(segments) != 2 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, forbiddenPayload)
		return
	}

	if strings.ToLower(segments[0]) != "bearer" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, forbiddenPayload)
		return
	}

	claims, err := h.infra.VerifyAndParseJWTToken(segments[1])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, forbiddenPayload)
		return
	}

	userIDf64, found := claims["user_id"].(float64)
	if !found {
		ctx.AbortWithStatusJSON(http.StatusForbidden, forbiddenPayload)
		return
	}

	userID := int(userIDf64)

	userInfo, err := h.usecase.GetUserInfo(ctx, userID)
	if err != nil {
		log.Error(ctx, userID, err, "failed to get user info")
		ctx.AbortWithStatusJSON(http.StatusForbidden, forbiddenPayload)
		return
	}

	ctx.Set(string(constants.CtxKeyUserID), userID)
	ctx.Set(string(constants.CtxKeyUserRole), userInfo.Role)

	ctx.Next()
}

func (h *handler) onlyAdmin(ctx *gin.Context) {
	userID := ctx.GetInt(string(constants.CtxKeyUserID))
	userInfo, err := h.usecase.GetUserInfo(ctx, userID)
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

	ctx.AbortWithStatusJSON(http.StatusForbidden, forbiddenPayload)
}
