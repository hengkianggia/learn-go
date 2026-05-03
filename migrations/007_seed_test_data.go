package migrations

import (
	"fmt"
	"learn/internal/config"
	"learn/internal/model"
	"learn/internal/pkg/qrcode"
	"learn/internal/pkg/random"
	"learn/internal/pkg/slug"
	"time"

	"gorm.io/gorm"
)

func init() {
	RegisterMigration("007", "Seed test data", migrate007)
}

func migrate007(db *gorm.DB) error {
	if isMigrationApplied(db, "007") {
		fmt.Println("Migration 007 already applied, skipping...")
		return nil
	}

	if err := seedTestGuests(db); err != nil {
		return err
	}
	if err := seedTestVenues(db); err != nil {
		return err
	}
	if err := seedTestUsers(db); err != nil {
		return err
	}
	eventIDs, priceIDs, err := seedTestEvents(db)
	if err != nil {
		return err
	}
	if err := seedTestOrders(db, eventIDs, priceIDs); err != nil {
		return err
	}

	recordMigration(db, "007", "Seed test data: 5 users, 3 guests, 3 venues, 3 events, orders, tickets + QR")
	return nil
}

func seedTestGuests(db *gorm.DB) error {
	guests := []model.Guest{
		{Name: "Dr. John Speaker", Slug: slug.GenerateSlug("Dr. John Speaker"), Bio: "Senior Software Engineer at Google"},
		{Name: "Jane Expert", Slug: slug.GenerateSlug("Jane Expert"), Bio: "Cloud Architect at AWS"},
		{Name: "Prof. Mike Talks", Slug: slug.GenerateSlug("Prof. Mike Talks"), Bio: "AI Researcher at OpenAI"},
	}
	for _, g := range guests {
		if err := db.Where("slug = ?", g.Slug).FirstOrCreate(&model.Guest{}, g).Error; err != nil {
			return fmt.Errorf("failed to seed guest %s: %w", g.Name, err)
		}
	}
	fmt.Println("3 guests created")
	return nil
}

func seedTestVenues(db *gorm.DB) error {
	venues := []model.Venue{
		{Name: "Jakarta Convention Center", Slug: slug.GenerateSlug("Jakarta Convention Center"), Address: "Jl. Gatot Subroto No.1, Jakarta", City: "Jakarta", State: "DKI Jakarta", Capacity: 5000, IsActive: true, Country: "Indonesia"},
		{Name: "Bali Nusa Dua Hall", Slug: slug.GenerateSlug("Bali Nusa Dua Hall"), Address: "Kawasan Nusa Dua, Bali", City: "Badung", State: "Bali", Capacity: 3000, IsActive: true, Country: "Indonesia"},
		{Name: "Bandung Creative Space", Slug: slug.GenerateSlug("Bandung Creative Space"), Address: "Jl. Braga No.50, Bandung", City: "Bandung", State: "Jawa Barat", Capacity: 500, IsActive: true, Country: "Indonesia"},
	}
	for _, v := range venues {
		if err := db.Where("slug = ?", v.Slug).FirstOrCreate(&model.Venue{}, v).Error; err != nil {
			return fmt.Errorf("failed to seed venue %s: %w", v.Name, err)
		}
	}
	fmt.Println("3 venues created")
	return nil
}

func seedTestUsers(db *gorm.DB) error {
	users := []model.User{
		{Name: "Alice Attendee", Email: "alice@test.com", Password: "password123", UserType: model.Attendee, IsVerified: true},
		{Name: "Bob Buyer", Email: "bob@test.com", Password: "password123", UserType: model.Attendee, IsVerified: true},
		{Name: "Charlie Customer", Email: "charlie@test.com", Password: "password123", UserType: model.Attendee, IsVerified: true},
		{Name: "Diana Devotee", Email: "diana@test.com", Password: "password123", UserType: model.Attendee, IsVerified: true},
		{Name: "Eve Enthusiast", Email: "eve@test.com", Password: "password123", UserType: model.Attendee, IsVerified: true},
	}
	for _, u := range users {
		if err := db.Where("email = ?", u.Email).FirstOrCreate(&model.User{}, u).Error; err != nil {
			return fmt.Errorf("failed to seed user %s: %w", u.Name, err)
		}
	}
	fmt.Println("5 users created")
	return nil
}

func seedTestEvents(db *gorm.DB) ([]uint, map[string]uint, error) {
	var jktVenue, baliVenue, bdgVenue model.Venue
	db.Where("slug = ?", slug.GenerateSlug("Jakarta Convention Center")).First(&jktVenue)
	db.Where("slug = ?", slug.GenerateSlug("Bali Nusa Dua Hall")).First(&baliVenue)
	db.Where("slug = ?", slug.GenerateSlug("Bandung Creative Space")).First(&bdgVenue)

	events := []struct {
		model.Event
		prices []struct{ name string; price int64; quota int }
	}{
		{
			Event: model.Event{
				Name: "Go Conference 2026", Slug: slug.GenerateSlug("Go Conference 2026"),
				Description: "The biggest Go conference in Southeast Asia", VenueID: jktVenue.ID,
				EventStartAt: time.Now().AddDate(0, 1, 0), Status: model.Published,
				SalesStartDate: time.Now().AddDate(0, 0, -7), SalesEndDate: time.Now().AddDate(0, 0, 25),
			},
			prices: []struct{ name string; price int64; quota int }{
				{"Early Bird", 150000, 200}, {"Regular", 250000, 500}, {"VIP", 500000, 50},
			},
		},
		{
			Event: model.Event{
				Name: "Cloud Native Summit", Slug: slug.GenerateSlug("Cloud Native Summit"),
				Description: "Master cloud-native technologies and Kubernetes", VenueID: baliVenue.ID,
				EventStartAt: time.Now().AddDate(0, 2, 0), Status: model.Published,
				SalesStartDate: time.Now().AddDate(0, 0, -3), SalesEndDate: time.Now().AddDate(0, 0, 50),
			},
			prices: []struct{ name string; price int64; quota int }{
				{"General Admission", 350000, 300}, {"Premium", 750000, 100},
			},
		},
		{
			Event: model.Event{
				Name: "AI & Data Workshop", Slug: slug.GenerateSlug("AI & Data Workshop"),
				Description: "Hands-on workshop on AI and data science", VenueID: bdgVenue.ID,
				EventStartAt: time.Now().AddDate(0, 3, 0), Status: model.Published,
				SalesStartDate: time.Now().AddDate(0, 0, -1), SalesEndDate: time.Now().AddDate(0, 0, 60),
			},
			prices: []struct{ name string; price int64; quota int }{
				{"Student", 100000, 100}, {"Professional", 300000, 150},
			},
		},
	}

	var eventIDs []uint
	priceIDs := make(map[string]uint)

	for _, e := range events {
		if err := db.Where("slug = ?", e.Slug).FirstOrCreate(&model.Event{}, e.Event).Error; err != nil {
			return nil, nil, fmt.Errorf("failed to seed event %s: %w", e.Name, err)
		}

		var savedEvent model.Event
		db.Where("slug = ?", e.Slug).First(&savedEvent)
		eventIDs = append(eventIDs, savedEvent.ID)

		for _, p := range e.prices {
			price := model.EventPrice{EventID: savedEvent.ID, Name: p.name, Price: p.price, Quota: p.quota}
			if err := db.Where("event_id = ? AND name = ?", savedEvent.ID, p.name).FirstOrCreate(&model.EventPrice{}, price).Error; err != nil {
				return nil, nil, fmt.Errorf("failed to seed event price %s: %w", p.name, err)
			}
			var savedPrice model.EventPrice
			db.Where("event_id = ? AND name = ?", savedEvent.ID, p.name).First(&savedPrice)
			priceIDs[savedEvent.Slug+"-"+p.name] = savedPrice.ID
		}
	}

	var guest1, guest2, guest3 model.Guest
	db.Where("slug = ?", slug.GenerateSlug("Dr. John Speaker")).First(&guest1)
	db.Where("slug = ?", slug.GenerateSlug("Jane Expert")).First(&guest2)
	db.Where("slug = ?", slug.GenerateSlug("Prof. Mike Talks")).First(&guest3)

	eventGuests := []model.EventGuest{
		{EventID: eventIDs[0], GuestID: guest1.ID, SessionTitle: "Building High-Performance Go Services"},
		{EventID: eventIDs[0], GuestID: guest2.ID, SessionTitle: "Go in Production at Scale"},
		{EventID: eventIDs[1], GuestID: guest2.ID, SessionTitle: "Kubernetes Best Practices"},
		{EventID: eventIDs[1], GuestID: guest3.ID, SessionTitle: "Serverless Architecture Patterns"},
		{EventID: eventIDs[2], GuestID: guest3.ID, SessionTitle: "Deep Learning with Go"},
	}
	for _, eg := range eventGuests {
		db.Where("event_id = ? AND guest_id = ?", eg.EventID, eg.GuestID).FirstOrCreate(&model.EventGuest{}, eg)
	}

	fmt.Println("3 events with prices and guest speakers created")
	return eventIDs, priceIDs, nil
}

func seedTestOrders(db *gorm.DB, eventIDs []uint, priceIDs map[string]uint) error {
	var users []model.User
	if err := db.Where("email IN ?", []string{"alice@test.com", "bob@test.com", "charlie@test.com", "diana@test.com", "eve@test.com"}).Find(&users).Error; err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	type orderSpec struct {
		userIdx      int
		eventSlug    string
		priceNames   []string
		quantities   []int
		seatNumbers  []string
	}
	orders := []orderSpec{
		{0, "go-conference-2026", []string{"Early Bird"}, []int{2}, []string{"A1", "A2"}},
		{1, "go-conference-2026", []string{"VIP"}, []int{1}, []string{"V1"}},
		{2, "cloud-native-summit", []string{"General Admission", "Premium"}, []int{1, 1}, []string{"B1", "P1"}},
		{3, "ai-data-workshop", []string{"Student"}, []int{1}, []string{"S1"}},
		{4, "ai-data-workshop", []string{"Professional"}, []int{2}, []string{"W1", "W2"}},
	}

	for _, spec := range orders {
		user := users[spec.userIdx]
		var totalPrice int64
		var lineItems []model.OrderLineItem

		for i, priceName := range spec.priceNames {
			priceID := priceIDs[spec.eventSlug+"-"+priceName]
			var eventPrice model.EventPrice
			db.First(&eventPrice, priceID)

			qty := spec.quantities[i]
			lineTotal := eventPrice.Price * int64(qty)
			totalPrice += lineTotal

			lineItems = append(lineItems, model.OrderLineItem{
				EventPriceID: priceID,
				Quantity:     qty,
				PricePerUnit: eventPrice.Price,
				TotalPrice:   lineTotal,
			})
		}

		order := model.Order{
			UserID:     user.ID,
			TotalPrice: totalPrice,
			Status:     model.OrderPaid,
			PaymentDue: time.Now().Add(24 * time.Hour),
		}

		if err := db.Create(&order).Error; err != nil {
			return fmt.Errorf("failed to create order for user %s: %w", user.Email, err)
		}

		for liIdx := range lineItems {
			lineItems[liIdx].OrderID = order.ID
		}
		if err := db.Create(&lineItems).Error; err != nil {
			return fmt.Errorf("failed to create order line items: %w", err)
		}

		var tickets []model.Ticket
		seatIdx := 0
		for _, li := range lineItems {
			for range li.Quantity {
				ticketCode := random.String(10)
				qrPath, qrErr := qrcode.GenerateQRCodePNG(config.AppConfig.StorageQRPath, ticketCode)
				if qrErr != nil {
					fmt.Printf("  [WARN] failed to generate QR for %s: %v\n", ticketCode, qrErr)
				}
				ticket := model.Ticket{
					OrderID:      order.ID,
					EventPriceID: li.EventPriceID,
					Price:        li.PricePerUnit,
					Type:         "General",
					TicketCode:   ticketCode,
					QrCodePath:   qrPath,
					SeatNumber:   spec.seatNumbers[seatIdx],
					OwnerName:    user.Name,
					OwnerEmail:   user.Email,
				}
				tickets = append(tickets, ticket)
				seatIdx++
			}
		}
		if err := db.Create(&tickets).Error; err != nil {
			return fmt.Errorf("failed to create tickets for order %d: %w", order.ID, err)
		}
		fmt.Printf("Order #%d for %s (%s) — %d tickets with QR codes\n", order.ID, user.Email, spec.eventSlug, len(tickets))
	}

	fmt.Println("Orders with PAID status, tickets, and QR codes generated")
	return nil
}

func isMigrationApplied(db *gorm.DB, version string) bool {
	var count int64
	db.Table("migrations").Where("version = ?", version).Count(&count)
	return count > 0
}

func recordMigration(db *gorm.DB, version, description string) {
	db.Table("migrations").Create(map[string]interface{}{
		"version":     version,
		"applied_at":  time.Now().Format(time.RFC3339),
		"description": description,
	})
}
