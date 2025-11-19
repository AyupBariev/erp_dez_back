package repositories

import (
	"erp/internal/app/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Получить заказы с фильтрацией по дате
func (r *OrderRepository) GetOrders(date *string) ([]*models.Order, error) {
	var orders []*models.Order

	query := withOrderPreloads(r.db.Model(&models.Order{}))

	if date != nil && *date != "" {
		query = query.Where("DATE(scheduled_at) = ?", *date)
	}

	if err := query.Order("scheduled_at DESC").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return orders, nil
}

// Получить следующий ERP-номер
func (r *OrderRepository) GetNextERPNumber() (int64, error) {
	var nextNumber int64
	err := r.db.Raw(`SELECT COALESCE(MAX(erp_number), 100000) + 1 FROM orders`).Scan(&nextNumber).Error
	return nextNumber, err
}

// Создание нового заказа
func (r *OrderRepository) Create(order *models.Order) error {
	if err := r.db.Create(order).Error; err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

// Получение заказа по первичному ключу
func (r *OrderRepository) GetOrderByID(orderID int64) (*models.Order, error) {
	var order models.Order
	err := withOrderPreloads(r.db).
		Where("id = ?", orderID).
		First(&order).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("order not found")
	}
	return &order, err
}

// Получение заказа по ERP номеру
func (r *OrderRepository) GetOrderByErpNumber(erpNumber int64) (*models.Order, error) {
	var order models.Order
	err := withOrderPreloads(r.db).
		Where("erp_number = ?", erpNumber).
		First(&order).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("order not found")
	}
	return &order, err
}

func (r *OrderRepository) GetByRepeatID(id uint) (*models.Order, error) {
	var order models.Order
	err := withOrderPreloads(r.db).
		Where("repeat_id = ?", id).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return &order, nil
}

// Заказы инженера за сегодня
func (r *OrderRepository) GetTodayOrders(chatID int64) ([]models.Order, error) {
	var orders []models.Order

	loc := time.Local
	nowLocal := time.Now().In(loc)

	startOfDayLocal := time.Date(
		nowLocal.Year(), nowLocal.Month(), nowLocal.Day(),
		0, 0, 0, 0, loc,
	)

	startOfDayUTC := startOfDayLocal.UTC()
	endOfDayUTC := startOfDayUTC.Add(24 * time.Hour)

	err := withOrderPreloads(r.db).
		Joins("JOIN engineers e ON e.id = orders.engineer_id").
		Where(`
            e.telegram_id = ?
            AND scheduled_at >= ?
            AND scheduled_at < ?
            AND status = ?
        `, chatID, startOfDayUTC, endOfDayUTC, "working").
		Find(&orders).Error

	return orders, err
}

// Повторные заказы инженера
func (r *OrderRepository) GetRepeatOrders(chatID int64) ([]models.Order, error) {
	var orders []models.Order

	err := withOrderPreloads(r.db).
		Joins("JOIN engineers e ON e.id = orders.engineer_id").
		Where("e.telegram_id = ? AND is_repeat = TRUE AND status = 'working'", chatID).
		Find(&orders).Error

	return orders, err
}

// Наличные заказы инженера
func (r *OrderRepository) GetCashOrders(chatID int64) ([]models.Order, error) {
	var orders []models.Order

	err := withOrderPreloads(r.db).
		Joins("JOIN engineers e ON e.id = orders.engineer_id").
		Where("e.telegram_id = ? AND status = 'sent_to_cash'", chatID).
		Find(&orders).Error

	return orders, err
}

func (r *OrderRepository) UpdateOrder(order *models.Order, fields ...string) error {
	if order.ID == 0 {
		return fmt.Errorf("order ID is required for update")
	}

	tx := r.db.Model(&models.Order{}).Where("id = ?", order.ID)
	if len(fields) > 0 {
		tx = tx.Select(fields)
	}

	return tx.Updates(order).Error
}

// --- вспомогательная функция ---
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

func (r *OrderRepository) Delete(erpNumber int64) error {
	return r.db.Delete(&models.Order{}, "erp_number = ?", erpNumber).Error
}
