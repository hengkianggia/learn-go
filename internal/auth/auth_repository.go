package auth

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user *User) error
	FindByUsername(username string) (*User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(user *User) error {
	// GORM BeforeSave hook akan hash password
	return r.db.Create(user).Error
}

func (r *userRepository) FindByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}