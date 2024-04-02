package usecase

import (
	"ara-server/internal/constants"
	"ara-server/internal/infrastructure/errors"
	"ara-server/internal/repository/db"
	"ara-server/util/log"
	"context"

	"golang.org/x/crypto/bcrypt"
)

var (
	errInvalidUser = errors.New("invalid user")
)

func (uc *Usecase) GetUserInfo(ctx context.Context, userID int) (UserInfo, error) {
	user, err := uc.db.GetUserInfoByUserID(ctx, userID)
	if err != nil {
		log.Error(ctx, userID, err, "failed to get user info")
		return UserInfo{}, err
	}

	if user.Name == "" {
		return UserInfo{}, errors.Wrap(errInvalidUser).WithCode(constants.ErrorCodeUserNotFound).WithType(errors.USER)
	}

	return UserInfo(user), nil
}

func (uc *Usecase) GetUserInfoList(ctx context.Context) ([]UserInfo, error) {
	users, err := uc.db.GetUserInfoList(ctx)
	if err != nil {
		log.Error(ctx, nil, err, "failed to get user info list")
		return nil, err
	}

	result := make([]UserInfo, 0, len(users))
	for _, user := range users {
		result = append(result, UserInfo(user))
	}

	return result, nil
}

func (uc *Usecase) LoginUser(ctx context.Context, param LoginUserParam) (AuthResponse, error) {
	user, err := uc.db.GetUserByEmail(ctx, param.Email)
	if err != nil {
		log.Error(ctx, param, err, "failed to get user by email")
		return AuthResponse{}, err
	}

	if user.ID == 0 {
		return AuthResponse{}, errors.Wrap(errInvalidUser).WithCode(constants.ErrorCodeUserNotFound).WithType(errors.USER)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password)); err != nil {
		return AuthResponse{}, errors.Wrap(errInvalidUser).WithType(errors.USER)
	}

	token, err := uc.infra.GenerateToken(map[string]interface{}{
		"user_id": user.ID,
	})
	if err != nil {
		log.Error(ctx, param, err, "failed to generate token")
		return AuthResponse{}, err
	}

	return AuthResponse{
		Token: token,
	}, nil
}

func (uc *Usecase) RegisterUser(ctx context.Context, param RegisterUserParam) (AuthResponse, error) {
	existingUser, err := uc.db.GetUserByEmail(ctx, param.Email)
	if err != nil {
		log.Error(ctx, param, err, "failed to get user by email")
		return AuthResponse{}, err
	}

	if existingUser.ID != 0 {
		return AuthResponse{}, errors.New("email already registered").WithCode(constants.ErrorCodeUserExists).WithType(errors.USER)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), 10)
	if err != nil {
		log.Error(ctx, param, err, "failed to hash password")
		return AuthResponse{}, err
	}

	tx, err := uc.db.BeginTx(ctx)
	if err != nil {
		log.Error(ctx, param, err, "failed to create tx")
		return AuthResponse{}, err
	}
	defer func() {
		if err == nil {
			return
		}

		if err := tx.Rollback(); err != nil {
			log.Error(ctx, param, err, "failed to rollback transaction")
		}
	}()

	userID, err := uc.db.InsertUser(ctx, tx, db.InsertUserParam{
		Email:    param.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		log.Error(ctx, param, err, "failed to insert user")
		return AuthResponse{}, err
	}

	err = uc.db.InsertProfile(ctx, tx, db.InsertProfileParam{
		UserID: userID,
		Name:   param.Name,
		Role:   constants.RoleUser,
	})
	if err != nil {
		log.Error(ctx, param, err, "failed to insert profile")
		return AuthResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		log.Error(ctx, param, err, "failed to commit transaction")
		return AuthResponse{}, err
	}

	token, err := uc.infra.GenerateToken(map[string]interface{}{
		"user_id": userID,
	})
	if err != nil {
		log.Error(ctx, param, err, "failed to generate token")
		return AuthResponse{}, err
	}

	return AuthResponse{
		Token: token,
	}, nil
}

func (uc *Usecase) UpdateUserInfo(ctx context.Context, userInfo UserInfo) error {
	user, err := uc.db.GetUserByID(ctx, userInfo.UserID)
	if err != nil {
		log.Error(ctx, userInfo, err, "failed to get user info")
		return err
	}

	if user.ID == 0 {
		return errors.Wrap(errInvalidUser).WithCode(constants.ErrorCodeUserNotFound).WithType(errors.USER)
	}

	tx, err := uc.db.BeginTx(ctx)
	if err != nil {
		log.Error(ctx, userInfo, err, "failed to create tx")
		return err
	}

	defer func() {
		if err == nil {
			return
		}

		if err := tx.Rollback(); err != nil {
			log.Error(ctx, userInfo, err, "failed to rollback transaction")
		}
	}()

	// update only if necessary
	if user.IsActive != userInfo.IsActive {
		user.IsActive = userInfo.IsActive
		err := uc.db.UpdateUser(ctx, tx, user)
		if err != nil {
			log.Error(ctx, userInfo, err, "failed to update user")
			return err
		}
	}

	if err := uc.db.UpdateUserProfile(ctx, tx, db.UserInfo{
		Name:   userInfo.Name,
		Role:   userInfo.Role,
		UserID: userInfo.UserID,
	}); err != nil {
		log.Error(ctx, userInfo, err, "failed to update user profile")
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error(ctx, userInfo, err, "failed to commit transaction")
		return err
	}

	return nil
}
