package auth

import (
	"errors"
	"learn/internal/model"
	"log/slog"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB, logger *slog.Logger) {
	users := []model.User{
		{
			Name:       "Administrator",
			Email:      "admin@example.com",
			Password:   "password123",
			UserType:   model.Administrator,
			IsVerified: true,
		},
		{
			Name:       "Organizer",
			Email:      "organizer@example.com",
			Password:   "password123",
			UserType:   model.Organizer,
			IsVerified: true,
		},
		{
			Name:       "hengki",
			Email:      "hengki@example.com",
			Password:   "password123",
			UserType:   model.Attendee,
			IsVerified: true,
		},
		{
			Name:       "anggia",
			Email:      "anggia@example.com",
			Password:   "password123",
			UserType:   model.Attendee,
			IsVerified: true,
		},
	}

	for _, user := range users {
		var existingUser model.User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&user).Error; err != nil {
					logger.Error("failed to seed user", slog.String("email", user.Email), slog.String("error", err.Error()))
				} else {
					logger.Info("seeded user", slog.String("email", user.Email))
				}
			}
		}
	}
}
