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
            address TEXT,
            city VARCHAR(255),
            country VARCHAR(255),
            capacity INTEGER
        )`,

		`CREATE TABLE IF NOT EXISTS guests (
            id SERIAL PRIMARY KEY,
            created_at TIMESTAMP WITH TIME ZONE,
            updated_at TIMESTAMP WITH TIME ZONE,
            deleted_at TIMESTAMP WITH TIME ZONE,
            name VARCHAR(255) NOT NULL,
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
            event_start_at TIMESTAMP WITH TIME ZONE,
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
            status VARCHAR(50) DEFAULT 'pending',
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
            type VARCHAR(255),
            ticket_code VARCHAR(255),
            owner_name VARCHAR(255),
            owner_email VARCHAR(255)
        )`,

		`CREATE TABLE IF NOT EXISTS payments (
            id SERIAL PRIMARY KEY,
            created_at TIMESTAMP WITH TIME ZONE,
            updated_at TIMESTAMP WITH TIME ZONE,
            deleted_at TIMESTAMP WITH TIME ZONE,
            order_id INTEGER NOT NULL,
            method VARCHAR(100),
            amount BIGINT,
            status VARCHAR(50),
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
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_user_id ON orders(user_id)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tickets_order_id ON tickets(order_id)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_payments_order_id ON payments(order_id)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_event_prices_event_id ON event_prices(event_id)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_events_deleted_at ON events(deleted_at)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_status ON orders(status)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_payments_status ON payments(status)`,
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

// Define the model structs for this migration
type User struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"not null"`
	Email          string `gorm:"unique;not null"`
	Password       string `gorm:"not null"`
	PhoneNumber    string
	ProfilePicture string
	UserType       string `gorm:"default:'attendee'"`
	IsVerified     bool   `gorm:"default:false"`
}

type Venue struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Address  string
	City     string
	Country  string
	Capacity int
}

type Guest struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Bio         string
	PhotoURL    string `gorm:"column:photo_url"`
	SocialMedia string
}

type Event struct {
	ID             uint   `gorm:"primaryKey"`
	VenueID        uint   `gorm:"not null"`
	Name           string `gorm:"not null"`
	Slug           string `gorm:"uniqueIndex;not null"`
	Description    string
	EventStartAt   string `gorm:"column:event_start_at"`
	Status         string `gorm:"default:'DRAFT'"`
	SalesStartDate string `gorm:"column:sales_start_date"`
	SalesEndDate   string `gorm:"column:sales_end_date"`
}

type EventPrice struct {
	ID      uint   `gorm:"primaryKey"`
	EventID uint   `gorm:"not null"`
	Name    string `gorm:"not null"`
	Price   int64  `gorm:"not null"`
	Quota   int    `gorm:"not null"`
}

type EventGuest struct {
	EventID      uint `gorm:"primaryKey"`
	GuestID      uint `gorm:"primaryKey"`
	SessionTitle string
}

type Order struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null"`
	TotalPrice  int64  `gorm:"not null"`
	Status      string `gorm:"default:'pending'"`
	PaymentDue  string `gorm:"column:payment_due"`
	ExpiredAt   string `gorm:"column:expired_at"`
	CancelledAt string `gorm:"column:cancelled_at"`
	CompletedAt string `gorm:"column:completed_at"`
}

type Ticket struct {
	ID           uint  `gorm:"primaryKey"`
	OrderID      uint  `gorm:"not null"`
	EventPriceID uint  `gorm:"not null"`
	Price        int64 `gorm:"not null"`
	Type         string
	TicketCode   string `gorm:"column:ticket_code"`
	OwnerName    string `gorm:"column:owner_name"`
	OwnerEmail   string `gorm:"column:owner_email"`
}

type Payment struct {
	ID      uint `gorm:"primaryKey"`
	OrderID uint `gorm:"not null"`
	Method  string
	Amount  int64
	Status  string
	Notes   string
}

type OrderLineItem struct {
	ID           uint  `gorm:"primaryKey"`
	OrderID      uint  `gorm:"not null"`
	EventPriceID uint  `gorm:"not null"`
	Quantity     int   `gorm:"not null"`
	PricePerUnit int64 `gorm:"not null"`
	TotalPrice   int64 `gorm:"not null"`
}
