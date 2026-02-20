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
	Name           string `gorm:"not null"`
	Email          string `gorm:"unique;not null"`
	Password       string `gorm:"not null"`
	PhoneNumber    string
	ProfilePicture string
	UserType       UserType `gorm:"default:'attendee'"`
	IsVerified     bool     `gorm:"default:false"`
	Orders         []Order
}

// BeforeSave is a GORM hook that hashes the user's password before saving.
// It skips hashing if the password is already a valid bcrypt hash to prevent double-hashing.
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	// bcrypt.Cost() returns nil error if the string is already a valid bcrypt hash
	if _, err := bcrypt.Cost([]byte(u.Password)); err == nil {
		return nil // already hashed, skip
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
