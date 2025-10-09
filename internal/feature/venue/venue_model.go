package venue

import (
	"gorm.io/gorm"
)

type Venue struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Address   string `gorm:"not null"`
	City      string
	State     string
	ZipCode   string
	Capacity  int
	IsActive  bool   `gorm:"default:true"`
}