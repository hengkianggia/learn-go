package migrations

import (
	"learn/internal/model"

	"gorm.io/gorm"
)

func init() {
	RegisterMigration("006", "Add QR code path to tickets", migrate006)
}

func migrate006(db *gorm.DB) error {
	return db.AutoMigrate(&model.Ticket{})
}
