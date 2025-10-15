package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserType string

const (
	Organizer     UserType = "organizer"
	Attendee      UserType = "attendee"
	Administrator UserType = "administrator"
)

type User struct {
	gorm.Model
	Name        string   `gorm:"not null"`
	Email       string   `gorm:"unique;not null"`
	Password    string   `gorm:"not null"`
	UserType    UserType `gorm:"type:user_type;not null"`
	PhoneNumber string
	IsVerified  bool `gorm:"default:false"`
}

// BeforeSave is a GORM hook that hashes the user's password before saving.
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
