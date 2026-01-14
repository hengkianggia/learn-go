package dto

import (
	"learn/internal/model"
	"time"
)

type TicketOrder struct {
	PriceId  string `json:"price_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1,max=10"`
}

type NewOrderInput struct {
	EventID        string        `json:"event_id" binding:"required"`
	TicketsOrdered []TicketOrder `json:"tickets_ordered" binding:"required,min=1,max=10,dive"`
}

type TicketResponse struct {
	ID         uint   `json:"id"`
	Price      int64  `json:"price"`      // Price in smallest currency unit (e.g., cents)
	Type       string `json:"type"`
	TicketCode string `json:"ticket_code"`
}

type OrderResponse struct {
	ID         uint              `json:"id"`
	TotalPrice int64             `json:"total_price"` // Total price in smallest currency unit (e.g., cents)
	Status     model.OrderStatus `json:"status"`
	PaymentDue time.Time         `json:"payment_due"`
	Tickets    []TicketResponse  `json:"tickets"`
}

func ToOrderResponse(order model.Order) OrderResponse {
	var ticketResponses []TicketResponse
	for _, ticket := range order.Tickets {
		ticketResponses = append(ticketResponses, TicketResponse{
			ID:         ticket.ID,
			Price:      ticket.Price,
			Type:       ticket.Type,
			TicketCode: ticket.TicketCode,
		})
	}

	return OrderResponse{
		ID:         order.ID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		PaymentDue: order.PaymentDue,
		Tickets:    ticketResponses,
	}
}
