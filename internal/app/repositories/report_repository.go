package repositories

import (
	"erp/internal/app/models"
	"erp/internal/domain"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) SaveReportLink(link *models.ReportLink) error {
	return r.db.Create(link).Error
}

func (r *ReportRepository) GetByToken(token string) (*models.ReportLink, error) {
	var link models.ReportLink
	if err := r.db.Where("token = ?", token).First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *ReportRepository) SaveReport(report *models.Report) error {
	return r.db.Save(report).Error
}

func (r *ReportRepository) Create(report *models.Report) error {
	return r.db.Create(report).Error
}

func (r *ReportRepository) Update(report *models.Report) error {
	return r.db.Save(report).Error
}

func (r *ReportRepository) GetByOrderID(orderID string) (*models.Report, error) {
	var report models.Report
	if err := r.db.Where("order_id = ?", orderID).First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) GetCashReports(from, to string) ([]*domain.CashReport, error) {
	fromTime, _ := time.Parse("2006-01-02", from)
	toTime, _ := time.Parse("2006-01-02", to)
	toTime = toTime.AddDate(0, 0, 1)

	var reports []*domain.CashReport
	err := r.db.
		Table("reports r").
		Select(`
			o.erp_number AS order_id,
			o.id AS id,
			o.engineer_id,
			CONCAT(COALESCE(e.first_name, ''), ' ', COALESCE(e.second_name, '')) as engineer_name,
			r.has_repeat, 
			r.repeat_date,
			r.repeat_note,
			o.note as description,
			o.scheduled_at,
			o.finish_price AS to_cash,
			COALESCE(r.gave_cash,0) AS gave_cash,
			r.issued_money,
			mm.motivation_percent,
			mm.total_motivation_amount * 0.5 AS prepayment,
			mm.total_motivation_amount * 0.5 AS salary,
			o.status as order_status
		`).
		Joins("LEFT JOIN orders o ON r.order_id = o.id").
		Joins("JOIN engineers e ON r.engineer_id = e.id").
		Joins("JOIN engineer_monthly_motivations mm ON mm.engineer_id = o.engineer_id AND DATE_FORMAT(o.scheduled_at,'%Y-%m') = DATE_FORMAT(mm.month,'%Y-%m')").
		Where("o.status in (?) AND o.scheduled_at BETWEEN ? AND ?", []string{"sent_to_cash", "closed_finally"}, fromTime, toTime).
		Order("o.scheduled_at DESC").
		Scan(&reports).Error

	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *ReportRepository) MarkIssued(orderID int64, gaveCash float64, issuedMoney float64, comment string) error {
	return r.db.Model(&models.Report{}).
		Where("order_id = ?", orderID).
		Updates(map[string]interface{}{
			"gave_cash":    gorm.Expr("gave_cash + ?", gaveCash),
			"issued_money": gorm.Expr("issued_money + ?", issuedMoney),
			"description":  gorm.Expr("CONCAT_WS('\n', description, ?)", comment),
		}).Error
}

// CheckMotivationPrepaymentLimit report_repository.go
// Проверяем право приёма денег по мотивационному ограничению.
// Возвращает ошибку, если ограничение НЕ пройдено.
func (r *ReportRepository) CheckMotivationPrepaymentLimit(orderID int64) error {
	type Info struct {
		MotivationPercent     float64
		TotalMotivationAmount float64
		PaidPrepayment        float64
	}

	var info Info
	err := r.db.Table("orders o").
		Select("mm.motivation_percent, mm.total_motivation_amount, COALESCE(SUM(ep.paid_prepayment), 0) AS paid_prepayment").
		Joins("JOIN engineer_monthly_motivations mm ON mm.engineer_id = o.engineer_id AND DATE_FORMAT(o.scheduled_at,'%Y-%m') = DATE_FORMAT(mm.month,'%Y-%m')").
		Joins("LEFT JOIN engineer_payouts ep ON ep.engineer_id = o.engineer_id").
		Where("o.id = ?", orderID).
		Group("mm.motivation_percent, mm.total_motivation_amount").
		Scan(&info).Error

	if err != nil {
		return fmt.Errorf("failed to check motivation limit: %w", err)
	}

	if info.MotivationPercent >= 20 &&
		info.PaidPrepayment <= info.TotalMotivationAmount*0.5 {
		return fmt.Errorf("prepayment limit exceeded: paid %.2f, allowed %.2f (50%% of %.2f)",
			info.PaidPrepayment, info.TotalMotivationAmount*0.5, info.TotalMotivationAmount)
	}

	return nil
}
