package migrations

import (
	"fmt"
	"learn/internal/model"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("004", "Seed initial data", migrate004)
}

func migrate004(db *gorm.DB) error {
	// Check if this migration was already applied
	var existingMigration MigrationHistory
	result := db.Where("version = ?", "004").First(&existingMigration)
	if result.Error == nil {
		fmt.Println("Migration 004 already applied, skipping...")
		return nil
	}

	// Check if admin user already exists
	var existingUser model.User
	result = db.Where("email = ?", "admin@example.com").First(&existingUser)
	if result.Error == nil {
		// Admin user already exists, skip seeding
		fmt.Println("Admin user already exists, skipping seeding...")
	} else {
		// Create admin user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		adminUser := model.User{
			Name:       "Administrator",
			Email:      "admin@example.com",
			Password:   string(hashedPassword),
			UserType:   "administrator",
			IsVerified: true,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		fmt.Println("Admin user created successfully")
	}

	// Record migration completion
	migration := MigrationHistory{
		Version:     "004",
		AppliedAt:   time.Now().Format(time.RFC3339),
		Description: "Seed initial data",
	}

	if err := db.Create(&migration).Error; err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	fmt.Println("Migration 004 completed successfully")
	return nil
}
