package router

import (
	"learn/internal/pkg/events"
	"learn/internal/repository"
	"log/slog"
)

// RegisterEventHandlersWithRepos registers all event handlers to the event bus with provided repositories
func RegisterEventHandlersWithRepos(eventBus *events.EventBus, orderRepo repository.OrderRepository,
	paymentRepo repository.PaymentRepository, ticketRepo repository.TicketRepository,
	eventRepo repository.EventRepository, logger *slog.Logger) {

	// Register OrderPaidEvent handler
	orderPaidHandler := events.NewOrderPaidEventHandler(orderRepo, logger)
	eventBus.Subscribe("order.paid", orderPaidHandler)

	// Register OrderCancelledEvent handler
	orderCancelledHandler := events.NewOrderCancelledEventHandler(orderRepo, logger)
	eventBus.Subscribe("order.cancelled", orderCancelledHandler)

	// Register PaymentStatusUpdatedEvent handler
	paymentStatusHandler := events.NewPaymentStatusUpdatedEventHandler(
		paymentRepo,
		orderRepo,
		ticketRepo,
		eventRepo,
		logger,
	)
	eventBus.Subscribe("payment.status.updated", paymentStatusHandler)
}
