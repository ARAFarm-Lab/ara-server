package usecase

import "ara-server/internal/constants"

type AuthResponse struct {
	Token string `json:"token,omitempty"`
}

type LoginUserParam struct {
	Email    string
	Password string
}

type RegisterUserParam struct {
	Email    string
	Password string
	Name     string
}

type UserInfo struct {
	UserID   int                `json:"user_id,omitempty"`
	IsActive bool               `json:"is_active"`
	Name     string             `json:"name"`
	Role     constants.UserRole `json:"role"`
}
