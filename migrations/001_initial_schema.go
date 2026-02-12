package migrations

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MigrationHistory struct {
	ID          uint   `gorm:"primaryKey"`
	Version     string `gorm:"not null"`
	AppliedAt   string `gorm:"not null"`
	Description string
}

func init() {
	RegisterMigration("001", "Create initial schema", migrate001)
}

func migrate001(db *gorm.DB) error {
	// Check if this migration was already applied
	var existingMigration MigrationHistory

	result := db.Where("version = ?", "001").First(&existingMigration)
	if result.Error == nil {
		fmt.Println("Migration 001 already applied, skipping...")
		return nil
	}

	// Create enums
	if err := createEnums(db); err != nil {
		return fmt.Errorf("failed to create enums: %w", err)
	}

	// Create tables using GORM AutoMigrate for the initial schema
	// This is safer than raw SQL for the initial setup
	if err := db.AutoMigrate(
		&MigrationHistory{},
	); err != nil {
		return fmt.Errorf("failed to create migration history table: %w", err)
	}

	// Record migration completion
	migration := MigrationHistory{
		Version:     "001",
		AppliedAt:   time.Now().Format(time.RFC3339),
		Description: "Create initial schema",
	}

	if err := db.Create(&migration).Error; err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	fmt.Println("Migration 001 completed successfully")
	return nil
}

func createEnums(db *gorm.DB) error {
	queries := []string{
		`DO $$ BEGIN
            CREATE TYPE user_type AS ENUM ('organizer', 'attendee', 'administrator');
        EXCEPTION
            WHEN duplicate_object THEN null;
        END $$;`,

		`DO $$ BEGIN
            CREATE TYPE event_status AS ENUM ('DRAFT', 'PUBLISHED', 'CANCELLED');
        EXCEPTION
            WHEN duplicate_object THEN null;
        END $$;`,
	}

	for _, query := range queries {
		if err := db.Exec(query).Error; err != nil {
			return err
		}
	}
	return nil
}
