package migrations

import (
	"learn/internal/model"

	"gorm.io/gorm"
)

func AddPaymentDetailsColumns(db *gorm.DB) error {
	// Add new columns to payments table
	// We use AutoMigrate with the updated model to add new columns
	err := db.AutoMigrate(&model.Payment{})
	if err != nil {
		return err
	}

	return nil
}
