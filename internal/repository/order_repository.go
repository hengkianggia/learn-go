package repository

import (
	"errors"
	"learn/internal/model"
	"time" // Added this import

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

type OrderRepository interface {
	CreateOrderInTransaction(order *model.Order, prices []model.EventPrice, priceUpdates map[uint]int) error
	GetEventPricesByIDs(priceIDs []uint) ([]model.EventPrice, error)
	GetEventByID(id uint) (*model.Event, error)
	GetOrderByID(orderID uint) (*model.Order, error)
	GetOrderByIDWithLineItems(orderID uint) (*model.Order, error) // Added
	UpdateOrder(order *model.Order) error
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrderInTransaction(order *model.Order, prices []model.EventPrice, priceUpdates map[uint]int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create Order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2. Create Payment record
		payment := &model.Payment{
			OrderID:       order.ID,
			PaymentMethod: model.PaymentMethodVirtualAccount, // Default or placeholder
			TransactionID: "",                                // Will be updated after payment gateway interaction
			PaymentStatus: model.PaymentStatusPending,
			PaymentDate:   time.Now(),
		}
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		// 3. Create OrderLineItems
		var orderLineItems []model.OrderLineItem
		priceMap := make(map[uint]model.EventPrice)
		for _, p := range prices {
			priceMap[p.ID] = p
		}

		for priceID, quantity := range priceUpdates {
			price, ok := priceMap[priceID]
			if !ok {
				return errors.New("price not found for order line item")
			}
			orderLineItems = append(orderLineItems, model.OrderLineItem{
				OrderID:      order.ID,
				EventPriceID: priceID,
				Quantity:     quantity,
				PricePerUnit: price.Price,
			})
		}
		if err := tx.Create(&orderLineItems).Error; err != nil {
			return err
		}

		// 4. Update quotas for each price
		for priceID, quantity := range priceUpdates {
			result := tx.Model(&model.EventPrice{}).Where("id = ? AND quota >= ? ", priceID, quantity).UpdateColumn("quota", gorm.Expr("quota - ?", quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("not enough quota or price not found")
			}
		}

		return nil
	})
}

func (r *orderRepository) GetEventPricesByIDs(priceIDs []uint) ([]model.EventPrice, error) {
	var prices []model.EventPrice
	err := r.db.Where("id IN ?", priceIDs).Find(&prices).Error
	return prices, err
}

func (r *orderRepository) GetEventByID(id uint) (*model.Event, error) {
	var event model.Event
	err := r.db.First(&event, id).Error
	return &event, err
}

func (r *orderRepository) GetOrderByID(orderID uint) (*model.Order, error) {
	var order model.Order
	if err := r.db.First(&order, orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateOrder(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) GetOrderByIDWithLineItems(orderID uint) (*model.Order, error) {
	var order model.Order
	if err := r.db.Preload("OrderLineItems").Preload("User").First(&order, orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
