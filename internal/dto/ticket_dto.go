package dto

import "learn/internal/model"

type CheckInTicketRequest struct {
	TicketCode string `json:"ticket_code" binding:"required"`
}

type CheckInTicketResponse struct {
	ID          uint   `json:"id"`
	TicketCode  string `json:"ticket_code"`
	Type        string `json:"type"`
	OwnerName   string `json:"owner_name"`
	OwnerEmail  string `json:"owner_email"`
	IsScanned   bool   `json:"is_scanned"`
	OrderID     uint   `json:"order_id"`
	OrderStatus string `json:"order_status"`
}

func ToCheckInTicketResponse(ticket model.Ticket, order model.Order) CheckInTicketResponse {
	return CheckInTicketResponse{
		ID:          ticket.ID,
		TicketCode:  ticket.TicketCode,
		Type:        ticket.Type,
		OwnerName:   ticket.OwnerName,
		OwnerEmail:  ticket.OwnerEmail,
		IsScanned:   ticket.IsScanned,
		OrderID:     ticket.OrderID,
		OrderStatus: string(order.Status),
	}
}
