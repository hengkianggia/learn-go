package dto

import (
	"learn/internal/model"
	"time"
)

type RegisterInput struct {
	Name            string         `json:"name" binding:"required"`
	Email           string         `json:"email" binding:"required,email"`
	Password        string         `json:"password" binding:"required"`
	ConfirmPassword string         `json:"confirm_password" binding:"required"`
	PhoneNumber     string         `json:"phone_number,omitempty"`
	UserType        model.UserType `json:"user_type,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerifyOTPInput struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type UserResponse struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phone_number,omitempty"`
	UserType    model.UserType `json:"user_type"`
	IsVerified  bool           `json:"is_verified"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
