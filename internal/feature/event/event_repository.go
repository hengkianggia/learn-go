package event

import "gorm.io/gorm"

type EventRepository interface {
	CreateEvent(event *Event) error
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
