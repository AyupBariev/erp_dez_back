package persistence

import (
	"errors"
	"fmt"
	"time"

	"erp/internal/domain/order"

	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{db: db}
}

// --- Получить все заказы (с фильтрацией по дате)
func (r *GormOrderRepository) GetOrders(date *string) ([]*order.Order, error) {
	var orders []*order.Order

	query := withOrderPreloads(r.db.Model(&order.Order{}))

	if date != nil && *date != "" {
		query = query.Where("DATE(scheduled_at) = ?", *date)
	}

	if err := query.Order("scheduled_at DESC").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return orders, nil
}

// --- Получить следующий ERP номер
func (r *GormOrderRepository) GetNextERPNumber() (int64, error) {
	var nextNumber int64
	err := r.db.Raw(`SELECT COALESCE(MAX(erp_number), 100000) + 1 FROM orders`).Scan(&nextNumber).Error
	return nextNumber, err
}

// --- Создать новый заказ
func (r *GormOrderRepository) Create(o *order.Order) error {
	if err := r.db.Create(o).Error; err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

// --- Получить заказ по ERP номеру
func (r *GormOrderRepository) GetByERPNumber(erpNumber int64) (*order.Order, error) {
	var o order.Order
	err := withOrderPreloads(r.db).
		Where("erp_number = ?", erpNumber).
		First(&o).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &o, err
}

// --- Заказы инженера за сегодня
func (r *GormOrderRepository) GetTodayOrders(chatID int64) ([]order.Order, error) {
	var orders []order.Order

	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := withOrderPreloads(r.db).
		Joins("JOIN engineers e ON e.id = orders.engineer_id").
		Where("e.telegram_id = ? AND scheduled_at BETWEEN ? AND ? AND status = ?", chatID, startOfDay, endOfDay, "working").
		Find(&orders).Error

	return orders, err
}

// --- Повторные заказы инженера
func (r *GormOrderRepository) GetRepeatOrders(chatID int64) ([]order.Order, error) {
	var orders []order.Order

	err := withOrderPreloads(r.db).
		Joins("JOIN engineers e ON e.id = orders.engineer_id").
		Where("e.telegram_id = ? AND is_repeat = TRUE AND status = 'working'", chatID).
		Find(&orders).Error

	return orders, err
}

// --- Наличные заказы инженера
func (r *GormOrderRepository) GetCashOrders(chatID int64) ([]order.Order, error) {
	var orders []order.Order

	err := withOrderPreloads(r.db).
		Joins("JOIN engineers e ON e.id = orders.engineer_id").
		Where("e.telegram_id = ? AND payment_type = 'cash'", chatID).
		Find(&orders).Error

	return orders, err
}

// --- Обновить заказ
func (r *GormOrderRepository) Update(o *order.Order) error {
	if o.ID == 0 {
		return fmt.Errorf("order ID is required for update")
	}

	return r.db.Model(&order.Order{}).
		Where("id = ?", o.ID).
		Updates(o).Error
}

// --- Удалить заказ
func (r *GormOrderRepository) Delete(erpNumber int64) error {
	return r.db.Delete(&order.Order{}, "erp_number = ?", erpNumber).Error
}

// --- Прелоады (ассоциации)
func withOrderPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Engineer").
		Preload("Aggregator", func(db *gorm.DB) *gorm.DB {
			return db.Table("aggregators")
		}).
		Preload("Problem", func(db *gorm.DB) *gorm.DB {
			return db.Table("problems")
		})
}
