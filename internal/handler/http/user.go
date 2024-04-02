package http

import (
	"ara-server/internal/constants"
	"ara-server/internal/infrastructure/errors"
	"ara-server/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetUserInfo(ctx *gin.Context) {
	userID := ctx.GetInt(string(constants.CtxKeyUserID))

	response, err := h.usecase.GetUserInfo(ctx, userID)
	if err != nil {
		if errors.IsUserError(err) {
			WriteJson(ctx, nil, err, http.StatusBadRequest)
			return
		}

		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, response, nil, http.StatusOK)
}

func (h *handler) HandleGetUserInfoList(ctx *gin.Context) {
	response, err := h.usecase.GetUserInfoList(ctx)
	if err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, response, nil)
}

func (h *handler) HandleLoginUser(ctx *gin.Context) {
	var request LoginUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		WriteJson(ctx, nil, errInvalidRequest, http.StatusBadRequest)
		return
	}

	response, err := h.usecase.LoginUser(ctx, usecase.LoginUserParam(request))
	if err != nil {
		if errors.IsUserError(err) {
			WriteJson(ctx, nil, err, http.StatusForbidden)
			return
		}

		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, response, nil, http.StatusOK)
}

func (h *handler) HandleRegisterUser(ctx *gin.Context) {
	var request RegisterUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		WriteJson(ctx, nil, errInvalidRequest, http.StatusBadRequest)
		return
	}

	response, err := h.usecase.RegisterUser(ctx, usecase.RegisterUserParam(request))
	if err != nil {
		if errors.IsUserError(err) {
			WriteJson(ctx, nil, err, http.StatusForbidden)
			return
		}

		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, response, nil, http.StatusCreated)
}

func (h *handler) HandleUpdateUserInfo(ctx *gin.Context) {
	var request UpdateUserInfoRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		WriteJson(ctx, nil, errInvalidRequest, http.StatusBadRequest)
		return
	}

	if err := h.usecase.UpdateUserInfo(ctx, usecase.UserInfo(request)); err != nil {
		WriteJson(ctx, nil, err)
		return
	}

	WriteJson(ctx, nil, nil, http.StatusAccepted)
}
