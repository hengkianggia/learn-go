package migration

import (
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Migration struct {
	Version string
	Name    string
	Func    MigrationFunc
}

type MigrationFunc func(*gorm.DB) error

var migrations = make(map[string]*Migration)

// RegisterMigration registers a new migration
func RegisterMigration(version, name string, fn MigrationFunc) {
	migrations[version] = &Migration{
		Version: version,
		Name:    name,
		Func:    fn,
	}
}

// Migrator handles database migrations
type Migrator struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB, logger *slog.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// Run runs all pending migrations
func (m *Migrator) Run() error {
	m.logger.Info("Starting migration process")

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Sort available migrations by version
	availableVersions := make([]string, 0, len(migrations))

	for version := range migrations {
		availableVersions = append(availableVersions, version)
	}

	sort.Slice(availableVersions, func(i, j int) bool {
		return compareVersions(availableVersions[i], availableVersions[j])
	})

	// Apply pending migrations
	for _, version := range availableVersions {
		if _, applied := appliedMigrations[version]; !applied {
			migration := migrations[version]

			m.logger.Info("Applying migration",
				slog.String("version", version),
				slog.String("name", migration.Name))

			if err := migration.Func(m.db); err != nil {
				m.logger.Error("Migration failed",
					slog.String("version", version),
					slog.String("error", err.Error()))
				return fmt.Errorf("migration %s failed: %w", version, err)
			}

			m.logger.Info("Migration completed", slog.String("version", version))
		}
	}

	m.logger.Info("All migrations completed")
	return nil
}

// GetAppliedMigrations returns a map of applied migrations
func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	var applied []struct {
		Version string
	}

	// Create migration history table if it doesn't exist
	if err := m.db.Table("migrations").AutoMigrate(&struct {
		ID          uint   `gorm:"primaryKey"`
		Version     string `gorm:"not null"`
		AppliedAt   string `gorm:"not null"`
		Description string
	}{}); err != nil {
		return nil, fmt.Errorf("failed to ensure migration history table: %w", err)
	}

	if err := m.db.Table("migrations").Select("version").Find(&applied).Error; err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	for _, migration := range applied {
		result[migration.Version] = true
	}
	return result, nil
}

// compareVersions compares two version strings numerically
func compareVersions(a, b string) bool {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		aNum, err1 := strconv.Atoi(aParts[i])
		bNum, err2 := strconv.Atoi(bParts[i])

		if err1 != nil || err2 != nil {
			// If conversion fails, fall back to string comparison
			return a < b
		}

		if aNum != bNum {
			return aNum < bNum
		}
	}

	return len(aParts) < len(bParts)
}
