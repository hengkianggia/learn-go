package dto

import (
	"learn/internal/model"
	"time"
)

type TicketOrder struct {
	PriceId  string `json:"price_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

type NewOrderInput struct {
	EventID        string        `json:"event_id" binding:"required"`
	TicketsOrdered []TicketOrder `json:"tickets_ordered" binding:"required"`
}

type TicketResponse struct {
	ID         uint    `json:"id"`
	Price      float64 `json:"price"`
	Type       string  `json:"type"`
	TicketCode string  `json:"ticket_code"`
}

type OrderResponse struct {
	ID         uint              `json:"id"`
	TotalPrice float64           `json:"total_price"`
	Status     model.OrderStatus `json:"status"`
	PaymentDue time.Time         `json:"payment_due"`
	Tickets    []TicketResponse  `json:"tickets"`
}

func ToOrderResponse(order model.Order) OrderResponse {
	var ticketResponses []TicketResponse
	for _, ticket := range order.Tickets {
		ticketResponses = append(ticketResponses, TicketResponse{
			ID:         ticket.ID,
			Price:      float64(ticket.Price),
			Type:       ticket.Type,
			TicketCode: ticket.TicketCode,
		})
	}

	return OrderResponse{
		ID:         order.ID,
		TotalPrice: float64(order.TotalPrice),
		Status:     order.Status,
		PaymentDue: order.PaymentDue,
		Tickets:    ticketResponses,
	}
}
