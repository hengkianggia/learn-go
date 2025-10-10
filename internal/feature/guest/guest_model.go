package guest

import (
	"gorm.io/gorm"
)

type Guest struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Slug      string `gorm:"uniqueIndex;not null"`
	Bio       string
}
