package migrations

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func init() {
	RegisterMigration("002", "Create user table", migrate002)
}

func migrate002(db *gorm.DB) error {
	// Check if this migration was already applied
	var existingMigration MigrationHistory

	result := db.Where("version = ?", "002").First(&existingMigration)

	if result.Error == nil {
		fmt.Println("Migration 002 already applied, skipping...")
		return nil
	}

	// Create users table
	query := `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        created_at TIMESTAMP WITH TIME ZONE,
        updated_at TIMESTAMP WITH TIME ZONE,
        deleted_at TIMESTAMP WITH TIME ZONE,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        phone_number VARCHAR(255),
        profile_picture VARCHAR(255),
        user_type user_type DEFAULT 'attendee',
        is_verified BOOLEAN DEFAULT FALSE
    )`

	if err := db.Exec(query).Error; err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Add indexes
	indexQueries := []string{
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users(email)`,
	}

	for _, query := range indexQueries {
		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	// Record migration completion
	migration := MigrationHistory{
		Version:     "002",
		AppliedAt:   time.Now().Format(time.RFC3339),
		Description: "Create user table",
	}

	if err := db.Create(&migration).Error; err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	fmt.Println("Migration 002 completed successfully")
	return nil
}
