package http

import "ara-server/internal/constants"

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInfoRequest struct {
	UserID   int                `json:"user_id"`
	IsActive bool               `json:"is_active,omitempty"`
	Name     string             `json:"name,omitempty"`
	Role     constants.UserRole `json:"role,omitempty"`
}
