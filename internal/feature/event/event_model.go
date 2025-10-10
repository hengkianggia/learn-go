package event

import (
	"learn/internal/feature/guest"
	"learn/internal/feature/venue"
	"time"

	"gorm.io/gorm"
)

type EventStatus string

const (
	Draft     EventStatus = "DRAFT"
	Published EventStatus = "PUBLISHED"
	Cancelled EventStatus = "CANCELLED"
)

type Event struct {
	gorm.Model
	VenueID        uint `gorm:"not null"`
	Venue          venue.Venue
	Name           string `gorm:"not null"`
	Description    string
	Date           time.Time   `gorm:"not null"`
	Time           time.Time   `gorm:"not null"`
	Status         EventStatus `gorm:"default:'DRAFT'"`
	SalesStartDate time.Time
	SalesEndDate   time.Time
	EventGuests    []EventGuest `gorm:"foreignKey:EventID"`
}

type EventGuest struct {
	EventID      uint `gorm:"primaryKey"`
	GuestID      uint `gorm:"primaryKey"`
	Event        Event
	Guest        guest.Guest
	SessionTitle string
}
