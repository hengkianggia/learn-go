package event

import "gorm.io/gorm"

type EventRepository interface {
	CreateEvent(event *Event) error
	GetEventByID(id uint) (*Event, error)
	CreateEventSpeakers(eventSpeakers []EventSpeaker) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) CreateEvent(event *Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) GetEventByID(id uint) (*Event, error) {
	var event Event
	err := r.db.Preload("Venue").Preload("EventSpeakers.Speaker").First(&event, id).Error
	return &event, err
}

func (r *eventRepository) CreateEventSpeakers(eventSpeakers []EventSpeaker) error {
	return r.db.Create(&eventSpeakers).Error
}
