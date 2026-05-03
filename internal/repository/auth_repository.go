package repository

import (
	"learn/internal/model"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	Save(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindAll(filters map[string]interface{}, page, limit int) ([]model.User, int64, error)
	Delete(id uint) error
	UpdateFields(id uint, fields map[string]interface{}) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(filters map[string]interface{}, page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
	for key, value := range filters {
		query = query.Where(key, value)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	if page <= 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 10
	}

	err := query.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}
