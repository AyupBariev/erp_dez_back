package repositories

import (
	"erp/internal/app/dto"
	"gorm.io/gorm"
	"time"
)

type EngineerMotivationRepository struct {
	DB *gorm.DB
}

func NewEngineerMotivationRepository(db *gorm.DB) *EngineerMotivationRepository {
	return &EngineerMotivationRepository{DB: db}
}

func (r *EngineerMotivationRepository) GetByDateRange(from, to time.Time) ([]dto.EngineerMotivationView, error) {
	var results []dto.EngineerMotivationView

	fromFormatted := from.Format("2006-01-02")
	toFormatted := to.Format("2006-01-02")

	err := r.DB.
		Table("engineers e").
		Select(`
			e.id as engineer_id,
			CONCAT(COALESCE(e.first_name, ''), ' ', COALESCE(e.second_name, '')) as engineer_name,
			m.month,
			COALESCE(m.reports_count, 0) as reports_count,
			COALESCE(m.primary_orders_count, 0) as primary_orders_count,
			COALESCE(m.repeat_orders_count, 0) as repeat_orders_count,
			COALESCE(m.orders_total_amount, 0) as orders_total_amount,
			COALESCE(m.repeat_orders_amount, 0) as repeat_orders_amount,
			COALESCE(m.gross_profit, 0) as gross_profit,
			COALESCE(m.average_check, 0) as average_check,
			COALESCE(m.total_motivation_amount, 0) as total_motivation_amount,
			COALESCE(m.aggregator_payout, 0) as aggregator_payout,
			COALESCE(m.gross_profit, 0) - COALESCE(m.total_motivation_amount, 0) as net_profit,
			COALESCE(m.gross_profit, 0) + COALESCE(m.aggregator_payout, 0) as total_amount,
			COALESCE(m.motivation_percent, 0) as motivation_percent
		`).
		Joins(`LEFT JOIN engineer_monthly_motivations m 
			   ON m.engineer_id = e.id 
			   AND m.month >= ? AND m.month <= ?`, fromFormatted, toFormatted).
		Group("e.id, e.first_name, e.second_name, m.month").
		Scan(&results).Error

	return results, err
}

// Получить все мотивации за месяц
func (r *EngineerMotivationRepository) GetByMonth(month time.Time) ([]dto.EngineerMotivationView, error) {
	var results []dto.EngineerMotivationView
	monthString := month.Format("2006-01") + "-01"
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
			COALESCE(m.gross_profit, 0) - COALESCE(m.total_motivation_amount, 0) as net_profit,
			COALESCE(m.average_check, 0) as average_check,
			COALESCE(m.motivation_percent, 0) as motivation_percent,
    		COALESCE(m.total_motivation_amount, 0) as total_motivation_amount,
    		COALESCE(m.aggregator_payout, 0) as aggregator_payout,
			COALESCE(m.gross_profit, 0) + COALESCE(m.aggregator_payout, 0) as total_amount
		`).
		Joins(`LEFT JOIN engineer_monthly_motivations m ON m.engineer_id = e.id AND m.month = ?`, monthString).
		Scan(&results).Error

	return results, err
}

func (r *EngineerMotivationRepository) GetProfitByDateRange(
	from, to time.Time,
) ([]dto.ProfitRaw, error) {

	var results []dto.ProfitRaw
	fromFormatted := from.Format("2006-01-02")
	toFormatted := to.Format("2006-01-02")
	err := r.DB.
		Table("engineer_monthly_motivations m").
		Select(`
			m.month AS date,
			SUM(
				COALESCE(m.gross_profit, 0)
				- COALESCE(m.total_motivation_amount, 0)
			) AS net_profit
		`).
		Where("m.month BETWEEN ? AND ?", fromFormatted, toFormatted).
		Group("m.month").
		Order("m.month").
		Scan(&results).Error

	return results, err
}
