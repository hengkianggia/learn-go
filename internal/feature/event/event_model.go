package event

import (
	"time"
	"learn/internal/feature/venue"
	"learn/internal/feature/speaker"
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
	VenueID       uint      `gorm:"not null"`
	Venue         venue.Venue
	Name          string    `gorm:"not null"`
	Description   string
	Date          time.Time `gorm:"not null"`
	Time          time.Time `gorm:"not null"`
	Status        EventStatus `gorm:"default:'DRAFT'"`
	SalesStartDate time.Time
	SalesEndDate  time.Time
}

type EventSpeaker struct {
	EventID      uint   `gorm:"primaryKey"`
	SpeakerID    uint   `gorm:"primaryKey"`
	Event        Event
	Speaker      speaker.Speaker
	SessionTitle string
}