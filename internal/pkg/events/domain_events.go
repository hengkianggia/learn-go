package events

import (
	"learn/internal/model"
	"time"
)

// OrderCreatedEvent is triggered when an order is created
type OrderCreatedEvent struct {
	OrderID     uint
	UserID      uint
	TotalPrice  int64
	CreatedAt   time.Time
}

func (e OrderCreatedEvent) GetEventType() string {
	return "order.created"
}

// OrderPaidEvent is triggered when an order is paid
type OrderPaidEvent struct {
	OrderID     uint
	UserID      uint
	TotalPrice  int64
	PaidAt      time.Time
}

func (e OrderPaidEvent) GetEventType() string {
	return "order.paid"
}

// OrderCancelledEvent is triggered when an order is cancelled
type OrderCancelledEvent struct {
	OrderID     uint
	UserID      uint
	Reason      string
	CancelledAt time.Time
}

func (e OrderCancelledEvent) GetEventType() string {
	return "order.cancelled"
}

// PaymentCreatedEvent is triggered when a payment is created
type PaymentCreatedEvent struct {
	PaymentID   uint
	OrderID     uint
	Method      model.PaymentMethod
	Amount      int64
	CreatedAt   time.Time
}

func (e PaymentCreatedEvent) GetEventType() string {
	return "payment.created"
}

// PaymentStatusUpdatedEvent is triggered when a payment status is updated
type PaymentStatusUpdatedEvent struct {
	PaymentID   uint
	OrderID     uint
	Status      model.PaymentStatus
	UpdatedAt   time.Time
}

func (e PaymentStatusUpdatedEvent) GetEventType() string {
	return "payment.status.updated"
}

// TicketsGeneratedEvent is triggered when tickets are generated for an order
type TicketsGeneratedEvent struct {
	OrderID     uint
	TicketCodes []string
	GeneratedAt time.Time
}

func (e TicketsGeneratedEvent) GetEventType() string {
	return "tickets.generated"
}