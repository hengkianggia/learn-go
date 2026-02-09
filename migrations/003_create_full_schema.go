package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

func init() {
	RegisterMigration("003", "Create full schema tables", migrate003)
}

func migrate003(db *gorm.DB) error {
	// Check if this migration was already applied
	var existingMigration MigrationHistory
	result := db.Where("version = ?", "003").First(&existingMigration)
	if result.Error == nil {
		fmt.Println("Migration 003 already applied, skipping...")
		return nil
	}

	// Create all tables using raw SQL for more control
	queries := []string{
		`CREATE TABLE IF NOT EXISTS venues (
            id SERIAL PRIMARY KEY,
            created_at TIMESTAMP WITH TIME ZONE,
            updated_at TIMESTAMP WITH TIME ZONE,
            deleted_at TIMESTAMP WITH TIME ZONE,
            name VARCHAR(255) NOT NULL,
            slug VARCHAR(255) UNIQUE NOT NULL,
            address TEXT NOT NULL,
            city VARCHAR(255),
            state VARCHAR(255),
            zip_code VARCHAR(255),
            country VARCHAR(255),
            capacity INTEGER,
            is_active BOOLEAN DEFAULT TRUE
        )`,

		`CREATE TABLE IF NOT EXISTS guests (
            id SERIAL PRIMARY KEY,
            created_at TIMESTAMP WITH TIME ZONE,
            updated_at TIMESTAMP WITH TIME ZONE,
            deleted_at TIMESTAMP WITH TIME ZONE,
            name VARCHAR(255) NOT NULL,
            slug VARCHAR(255) UNIQUE NOT NULL,
            bio TEXT,
            photo_url VARCHAR(255),
            social_media TEXT
        )`,

		`CREATE TABLE IF NOT EXISTS events (
            id SERIAL PRIMARY KEY,
            created_at TIMESTAMP WITH TIME ZONE,
            updated_at TIMESTAMP WITH TIME ZONE,
            deleted_at TIMESTAMP WITH TIME ZONE,
            venue_id INTEGER NOT NULL,
            name VARCHAR(255) NOT NULL,
            slug VARCHAR(255) UNIQUE,
            description TEXT,
            event_start_at TIMESTAMP WITH TIME ZONE NOT NULL,
            status event_status DEFAULT 'DRAFT',
            sales_start_date TIMESTAMP WITH TIME ZONE,
            sales_end_date TIMESTAMP WITH TIME ZONE
        )`,

		`CREATE TABLE IF NOT EXISTS event_prices (
            id SERIAL PRIMARY KEY,
            created_at TIMESTAMP WITH TIME ZONE,
            updated_at TIMESTAMP WITH TIME ZONE,
            deleted_at TIMESTAMP WITH TIME ZONE,
            event_id INTEGER NOT NULL,
            name VARCHAR(255) NOT NULL,
            price BIGINT NOT NULL,
            quota INTEGER NOT NULL
        )`,

		`CREATE TABLE IF NOT EXISTS event_guests (
            event_id INTEGER NOT NULL,
            guest_id INTEGER NOT NULL,
            session_title VARCHAR(255),
            PRIMARY KEY (event_id, guest_id)
        )`,

		`CREATE TABLE IF NOT EXISTS orders (
            id SERIAL PRIMARY KEY,
            created_at TIMESTAMP WITH TIME ZONE,
            updated_at TIMESTAMP WITH TIME ZONE,
            deleted_at TIMESTAMP WITH TIME ZONE,
            user_id INTEGER NOT NULL,
            total_price BIGINT NOT NULL,
            status VARCHAR(50) DEFAULT 'PENDING' NOT NULL,
            payment_due TIMESTAMP WITH TIME ZONE,
            expired_at TIMESTAMP WITH TIME ZONE,
            cancelled_at TIMESTAMP WITH TIME ZONE,
            completed_at TIMESTAMP WITH TIME ZONE
        )`,

		`CREATE TABLE IF NOT EXISTS tickets (
            id SERIAL PRIMARY KEY,
            created_at TIMESTAMP WITH TIME ZONE,
            updated_at TIMESTAMP WITH TIME ZONE,
            deleted_at TIMESTAMP WITH TIME ZONE,
            order_id INTEGER NOT NULL,
            event_price_id INTEGER NOT NULL,
            price BIGINT NOT NULL,
            type VARCHAR(255) NOT NULL,
            ticket_code VARCHAR(255) UNIQUE NOT NULL,
            seat_number VARCHAR(255),
            is_scanned BOOLEAN DEFAULT FALSE,
            owner_name VARCHAR(255),
            owner_email VARCHAR(255)
        )`,

		`CREATE TABLE IF NOT EXISTS payments (
            id SERIAL PRIMARY KEY,
            created_at TIMESTAMP WITH TIME ZONE,
            updated_at TIMESTAMP WITH TIME ZONE,
            deleted_at TIMESTAMP WITH TIME ZONE,
            order_id INTEGER NOT NULL,
            payment_method VARCHAR(50) NOT NULL,
            transaction_id VARCHAR(255) UNIQUE NOT NULL,
            amount BIGINT,
            payment_status VARCHAR(50) NOT NULL,
            payment_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            notes TEXT
        )`,

		`CREATE TABLE IF NOT EXISTS order_line_items (
            id SERIAL PRIMARY KEY,
            created_at TIMESTAMP WITH TIME ZONE,
            updated_at TIMESTAMP WITH TIME ZONE,
            deleted_at TIMESTAMP WITH TIME ZONE,
            order_id INTEGER NOT NULL,
            event_price_id INTEGER NOT NULL,
            quantity INTEGER NOT NULL,
            price_per_unit BIGINT NOT NULL,
            total_price BIGINT NOT NULL
        )`,
	}

	for _, query := range queries {
		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("failed to execute migration query: %w", err)
		}
	}

	// Add indexes
	indexQueries := []string{
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_events_slug ON events(slug)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_venues_slug ON venues(slug)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_guests_slug ON guests(slug)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_user_id ON orders(user_id)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tickets_order_id ON tickets(order_id)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_payments_order_id ON payments(order_id)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_event_prices_event_id ON event_prices(event_id)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_events_deleted_at ON events(deleted_at)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_status ON orders(status)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_payments_payment_status ON payments(payment_status)`,
	}

	for _, query := range indexQueries {
		if err := db.Exec(query).Error; err != nil {
			// Continue even if index creation fails (might already exist)
			fmt.Printf("Warning: failed to create index: %v\n", err)
		}
	}

	// Record migration completion
	migration := MigrationHistory{
		Version:     "003",
		AppliedAt:   "NOW()",
		Description: "Create full schema tables",
	}

	if err := db.Create(&migration).Error; err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	fmt.Println("Migration 003 completed successfully")
	return nil
}
