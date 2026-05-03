package migrations

import (
	"learn/internal/model"

	"gorm.io/gorm"
)

func init() {
	RegisterMigration("008", "add_user_approval_blocked", AddUserApprovalBlockedColumns)
}

func AddUserApprovalBlockedColumns(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return err
	}

	return db.Exec("UPDATE users SET is_approved = true, is_blocked = false WHERE is_approved = false AND is_blocked = false").Error
}
