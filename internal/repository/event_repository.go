package repository

import (
	"learn/internal/model"

	"gorm.io/gorm"
)

type EventRepository interface {
	CreateEvent(event *model.Event) error
	GetEventByID(id uint) (*model.Event, error)
	FindBySlug(slug string) (*model.Event, error)
	CreateEventGuests(eventGuests []model.EventGuest) error
	CreateEventPrices(eventPrices []model.EventPrice) error
	GetEventPriceByID(id uint) (*model.EventPrice, error) // Added
	GetEventsByGuestSlug(guestSlug string) ([]model.Event, error)
	UpdateEvent(event *model.Event) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) CreateEvent(event *model.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) GetEventByID(id uint) (*model.Event, error) {
	var event model.Event
	err := r.db.Preload("Venue").Preload("EventGuests.Guest").Preload("Prices").First(&event, id).Error
	return &event, err
}

func (r *eventRepository) FindBySlug(slug string) (*model.Event, error) {
	var event model.Event
	err := r.db.Preload("Venue").Preload("EventGuests.Guest").Preload("Prices").Where("slug = ?", slug).First(&event).Error
	return &event, err
}

func (r *eventRepository) CreateEventGuests(eventGuests []model.EventGuest) error {
	return r.db.Create(&eventGuests).Error
}

func (r *eventRepository) CreateEventPrices(eventPrices []model.EventPrice) error {
	return r.db.Create(&eventPrices).Error
}

func (r *eventRepository) GetEventsByGuestSlug(guestSlug string) ([]model.Event, error) {
	var events []model.Event
	err := r.db.Joins("JOIN event_guests ON event_guests.event_id = events.id").Joins("JOIN guests ON guests.id = event_guests.guest_id").Where("guests.slug = ?", guestSlug).Preload("Venue").Preload("EventGuests.Guest").Preload("Prices").Find(&events).Error
	return events, err
}

func (r *eventRepository) UpdateEvent(event *model.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) GetEventPriceByID(id uint) (*model.EventPrice, error) {
	var eventPrice model.EventPrice
	err := r.db.First(&eventPrice, id).Error
	return &eventPrice, err
}
