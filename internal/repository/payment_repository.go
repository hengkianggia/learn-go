package repository

import (
	"learn/internal/config"
	"learn/internal/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(payment *model.Payment) error
	GetPaymentByID(paymentID uint) (*model.Payment, error)
	GetPaymentByOrderID(orderID uint) (*model.Payment, error)
	UpdatePayment(payment *model.Payment) error
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

func (r *paymentRepository) DeletePayment(paymentID uint) error {
	return r.db.Delete(&model.Payment{}, paymentID).Error
}
