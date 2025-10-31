package repositories

import (
	"erp/internal/app/models"
	"gorm.io/gorm"
	"time"
)

type EngineerMotivationRepository struct {
	DB *gorm.DB
}

func NewEngineerMotivationRepository(db *gorm.DB) *EngineerMotivationRepository {
	return &EngineerMotivationRepository{DB: db}
}

// Получить все мотивации за месяц
func (r *EngineerMotivationRepository) GetByMonth(month time.Time) ([]models.EngineerMonthlyMotivation, error) {
	var results []models.EngineerMonthlyMotivation

	err := r.DB.
		Table("engineers e").
		Select(`
			e.id as engineer_id,
    		CONCAT(COALESCE(e.first_name, ''), ' ', COALESCE(e.second_name, '')) as engineer_name,
			COALESCE(m.reports_count, 0) as reports_count,
			COALESCE(m.primary_orders_count, 0) as primary_orders_count,
			COALESCE(m.repeat_orders_count, 0) as repeat_orders_count,
			COALESCE(m.orders_total_amount, 0) as orders_total_amount,
			COALESCE(m.repeat_orders_amount, 0) as repeat_orders_amount,
			COALESCE(m.gross_profit, 0) as gross_profit,
			COALESCE(m.average_check, 0) as average_check,
			COALESCE(m.motivation_percent, 0) as motivation_percent,
			COALESCE(m.total_motivation, 0) as total_motivation,
			COALESCE(m.confirmed_by_admin, false) as confirmed_by_admin
		`).
		Joins(`LEFT JOIN engineer_monthly_motivations m ON m.engineer_id = e.id AND m.month = ?`, month).
		Scan(&results).Error

	return results, err
}
