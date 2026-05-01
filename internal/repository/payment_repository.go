package repository

import (
	"learn/internal/config"
	"learn/internal/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentRepository interface {
	CreatePayment(payment *model.Payment) error
	CreatePaymentInTransaction(payment *model.Payment) error
	GetPaymentByID(paymentID uint) (*model.Payment, error)
	GetPaymentByOrderID(orderID uint) (*model.Payment, error)
	GetPaymentByTransactionID(transactionID string) (*model.Payment, error)
	UpdatePayment(payment *model.Payment) error
	UpdatePaymentStatusInTransaction(paymentID uint, status model.PaymentStatus) (*model.Payment, bool, error)
	DeletePayment(paymentID uint) error
	GetRedisClient() *redis.Client
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) GetRedisClient() *redis.Client {
	return config.Rdb
}

func (r *paymentRepository) CreatePayment(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) CreatePaymentInTransaction(payment *model.Payment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, payment.OrderID).Error; err != nil {
			return err
		}
		if order.Status != model.OrderPending {
			return gorm.ErrInvalidTransaction
		}

		var existing model.Payment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ?", payment.OrderID).First(&existing).Error; err == nil {
			return gorm.ErrDuplicatedKey
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		return tx.Create(payment).Error
	})
}

func (r *paymentRepository) GetPaymentByID(paymentID uint) (*model.Payment, error) {
	var payment model.Payment
	if err := r.db.First(&payment, paymentID).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetPaymentByOrderID(orderID uint) (*model.Payment, error) {
	var payment model.Payment
	if err := r.db.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdatePayment(payment *model.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) UpdatePaymentStatusInTransaction(paymentID uint, status model.PaymentStatus) (*model.Payment, bool, error) {
	var payment model.Payment
	changed := false

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&payment, paymentID).Error; err != nil {
			return err
		}

		if payment.PaymentStatus == status {
			return nil
		}

		payment.PaymentStatus = status
		if err := tx.Save(&payment).Error; err != nil {
			return err
		}

		switch status {
		case model.PaymentStatusSuccess:
			if err := tx.Model(&model.Order{}).Where("id = ?", payment.OrderID).Update("status", model.OrderPaid).Error; err != nil {
				return err
			}
		case model.PaymentStatusFailed:
			if err := tx.Model(&model.Order{}).Where("id = ? AND status = ?", payment.OrderID, model.OrderPending).Update("status", model.OrderCancelled).Error; err != nil {
				return err
			}
		}

		changed = true
		return nil
	})
	if err != nil {
		return nil, false, err
	}
	return &payment, changed, nil
}

func (r *paymentRepository) DeletePayment(paymentID uint) error {
	return r.db.Delete(&model.Payment{}, paymentID).Error
}

func (r *paymentRepository) GetPaymentByTransactionID(transactionID string) (*model.Payment, error) {
	var payment model.Payment
	if err := r.db.Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}
