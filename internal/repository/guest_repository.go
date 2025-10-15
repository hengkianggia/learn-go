package repository

import (
	"learn/internal/model"

	"gorm.io/gorm"
)

type GuestRepository interface {
	CreateGuest(guest *model.Guest) error
	GetGuestByID(id uint) (*model.Guest, error)
	FindBySlug(slug string) (*model.Guest, error)
}

type guestRepository struct {
	db *gorm.DB
}

func NewGuestRepository(db *gorm.DB) GuestRepository {
	return &guestRepository{db: db}
}

func (r *guestRepository) CreateGuest(guest *model.Guest) error {
	return r.db.Create(guest).Error
}

func (r *guestRepository) GetGuestByID(id uint) (*model.Guest, error) {
	var guest model.Guest
	err := r.db.First(&guest, id).Error
	return &guest, err
}

func (r *guestRepository) FindBySlug(slug string) (*model.Guest, error) {
	var guest model.Guest
	err := r.db.Where("slug = ?", slug).First(&guest).Error
	return &guest, err
}
