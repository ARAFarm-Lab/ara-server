package http

import (
	"ara-server/internal/constants"
	"ara-server/internal/infrastructure/errors"
	"ara-server/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetUserInfo(ctx *gin.Context) error {
	userID := ctx.GetInt(string(constants.CtxKeyUserID))

	response, err := h.usecase.GetUserInfo(ctx, userID)
	if err != nil {
		if errors.IsUserError(err) {
			WriteJson(ctx, nil, err, http.StatusBadRequest)
			return err
		}

		WriteJson(ctx, nil, err)
		return err
	}

	WriteJson(ctx, response, nil, http.StatusOK)
	return nil
}

func (h *handler) HandleGetUserInfoList(ctx *gin.Context) error {
	response, err := h.usecase.GetUserInfoList(ctx)
	if err != nil {
		WriteJson(ctx, nil, err)
		return err
	}

	WriteJson(ctx, response, nil)
	return nil
}

func (h *handler) HandleLoginUser(ctx *gin.Context) error {
	var request LoginUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		WriteJson(ctx, nil, errInvalidRequest, http.StatusBadRequest)
		return errInvalidRequest
	}

	response, err := h.usecase.LoginUser(ctx, usecase.LoginUserParam(request))
	if err != nil {
		if errors.IsUserError(err) {
			WriteJson(ctx, nil, err, http.StatusForbidden)
			return err
		}

		WriteJson(ctx, nil, err)
		return err
	}

	WriteJson(ctx, response, nil, http.StatusOK)
	return nil
}

func (h *handler) HandleRegisterUser(ctx *gin.Context) error {
	var request RegisterUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		WriteJson(ctx, nil, errInvalidRequest, http.StatusBadRequest)
		return errInvalidRequest
	}

	response, err := h.usecase.RegisterUser(ctx, usecase.RegisterUserParam(request))
	if err != nil {
		if errors.IsUserError(err) {
			WriteJson(ctx, nil, err, http.StatusForbidden)
			return err
		}

		WriteJson(ctx, nil, err)
		return err
	}

	WriteJson(ctx, response, nil, http.StatusCreated)
	return nil
}

func (h *handler) HandleUpdateUserInfo(ctx *gin.Context) error {
	var request UpdateUserInfoRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		WriteJson(ctx, nil, errInvalidRequest, http.StatusBadRequest)
		return errInvalidRequest
	}

	if err := h.usecase.UpdateUserInfo(ctx, usecase.UserInfo(request)); err != nil {
		WriteJson(ctx, nil, err)
		return err
	}

	WriteJson(ctx, nil, nil, http.StatusAccepted)
	return nil
}
