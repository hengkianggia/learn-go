package repository

import (
	"errors"
	"learn/internal/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderInTransaction(order *model.Order, tickets []model.Ticket, priceUpdates map[uint]int) error
	GetEventPricesByIDs(priceIDs []uint) ([]model.EventPrice, error)
	GetEventByID(id uint) (*model.Event, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrderInTransaction(order *model.Order, tickets []model.Ticket, priceUpdates map[uint]int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create Order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2. Assign OrderID to tickets
		for i := range tickets {
			tickets[i].OrderID = order.ID
		}

		// 3. Bulk create tickets
		if err := tx.Create(&tickets).Error; err != nil {
			return err
		}

		// 4. Update quotas for each price
		for priceID, quantity := range priceUpdates {
			result := tx.Model(&model.EventPrice{}).Where("id = ? AND quota >= ?", priceID, quantity).UpdateColumn("quota", gorm.Expr("quota - ?", quantity))
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
