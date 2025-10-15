package repository

import (
	"learn/internal/model"

	"gorm.io/gorm"
)

type VenueRepository interface {
	CreateVenue(venue *model.Venue) error
	GetVenueByID(id uint) (*model.Venue, error)
	FindBySlug(slug string) (*model.Venue, error)
}

type venueRepository struct {
	db *gorm.DB
}

func NewVenueRepository(db *gorm.DB) VenueRepository {
	return &venueRepository{db: db}
}

func (r *venueRepository) CreateVenue(venue *model.Venue) error {
	return r.db.Create(venue).Error
}

func (r *venueRepository) GetVenueByID(id uint) (*model.Venue, error) {
	var venue model.Venue
	err := r.db.First(&venue, id).Error
	return &venue, err
}

func (r *venueRepository) FindBySlug(slug string) (*model.Venue, error) {
	var venue model.Venue
	err := r.db.Where("slug = ?", slug).First(&venue).Error
	return &venue, err
}
