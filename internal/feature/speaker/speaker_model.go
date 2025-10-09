package speaker

import (
	"gorm.io/gorm"
)

type Speaker struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique"`
	Bio       string
}