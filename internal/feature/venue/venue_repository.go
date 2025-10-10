package venue

import "gorm.io/gorm"

type VenueRepository interface {
	CreateVenue(venue *Venue) error
	GetVenueByID(id uint) (*Venue, error)
	FindBySlug(slug string) (*Venue, error)
}

type venueRepository struct {
	db *gorm.DB
}

func NewVenueRepository(db *gorm.DB) VenueRepository {
	return &venueRepository{db: db}
}

func (r *venueRepository) CreateVenue(venue *Venue) error {
	return r.db.Create(venue).Error
}

func (r *venueRepository) GetVenueByID(id uint) (*Venue, error) {
	var venue Venue
	err := r.db.First(&venue, id).Error
	return &venue, err
}

func (r *venueRepository) FindBySlug(slug string) (*Venue, error) {
	var venue Venue
	err := r.db.Where("slug = ?", slug).First(&venue).Error
	return &venue, err
}
