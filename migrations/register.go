package migrations

import (
	"learn/internal/migration"

	"gorm.io/gorm"
)

// RegisterMigration is a wrapper around the internal migration.RegisterMigration
func RegisterMigration(version, name string, fn func(*gorm.DB) error) {
	migration.RegisterMigration(version, name, fn)
}
