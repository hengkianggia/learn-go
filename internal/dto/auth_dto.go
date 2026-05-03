package dto

import (
	"learn/internal/model"
	"time"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerifyOTPInput struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type AdminUserActionInput struct {
	UserID uint `json:"user_id" binding:"required"`
}

type AdminUserListInput struct {
	UserType model.UserType `json:"user_type,omitempty" form:"user_type"`
	Page     int            `json:"page,omitempty" form:"page"`
	Limit    int            `json:"limit,omitempty" form:"limit"`
}

type RegisterInput struct {
	Name            string         `json:"name" binding:"required"`
	Email           string         `json:"email" binding:"required,email"`
	Password        string         `json:"password" binding:"required,min=6"`
	ConfirmPassword string         `json:"confirm_password" binding:"required"`
	PhoneNumber     string         `json:"phone_number,omitempty"`
	UserType        model.UserType `json:"user_type,omitempty"`
}

type UserResponse struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phone_number,omitempty"`
	UserType    model.UserType `json:"user_type"`
	IsVerified  bool           `json:"is_verified"`
	IsApproved  bool           `json:"is_approved"`
	IsBlocked   bool           `json:"is_blocked"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
