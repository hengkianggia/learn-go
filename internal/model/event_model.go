package model

import (
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
	Venue          Venue
	Name           string `gorm:"not null"`
	Slug           string `gorm:"uniqueIndex;not null"`
	Description    string
	EventStartAt   time.Time   `gorm:"not null"`
	Status         EventStatus `gorm:"default:'DRAFT'"`
	SalesStartDate time.Time
	SalesEndDate   time.Time
	EventGuests    []EventGuest `gorm:"foreignKey:EventID"`
	Prices         []EventPrice `gorm:"foreignKey:EventID"`
}

type EventGuest struct {
	EventID      uint `gorm:"primaryKey"`
	GuestID      uint `gorm:"primaryKey"`
	Event        Event
	Guest        Guest
	SessionTitle string
}

type EventPrice struct {
	gorm.Model
	EventID uint   `gorm:"not null"`
	Name    string `gorm:"not null"` // e.g., "Presale", "VIP"
	Price   int64  `gorm:"not null"` // Price in smallest currency unit (e.g., cents)
	Quota   int    `gorm:"not null"`
	Tickets []Ticket
}
