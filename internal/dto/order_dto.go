package dto

type TicketOrder struct {
	PriceId  string `json:"price_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

type NewOrderInput struct {
	// UserID         string        `json:"user_id" binding:"required"`
	EventID        string        `json:"event_id" binding:"required"`
	TicketsOrdered []TicketOrder `json:"tickets_ordered" binding:"required"`
}
