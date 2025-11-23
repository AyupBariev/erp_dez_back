package repositories

import (
	"database/sql"
	"erp/internal/app/models"
	"erp/internal/domain"
	"erp/internal/pkg/logger"
	"fmt"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) SaveReportLink(link *models.ReportLink) error {
	_, err := r.db.Exec(`
		INSERT INTO report_links (token, order_id, engineer_id, expires_at, created_at)
		VALUES (?, ?, ?, ?, NOW())`,
		link.Token, link.OrderID, link.EngineerID, link.ExpiresAt,
	)

	if err != nil {
		logger.LogError("Ошибка сохранения report_link.", err)
	}
	return err
}

func (r *ReportRepository) GetByToken(token string) (*models.ReportLink, error) {
	row := r.db.QueryRow(`
		SELECT id, order_id, engineer_id, token, expires_at, created_at
		FROM report_links
		WHERE token = ?`, token,
	)

	var link models.ReportLink
	if err := row.Scan(&link.ID, &link.OrderID, &link.EngineerID, &link.Token, &link.ExpiresAt, &link.CreatedAt); err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *ReportRepository) SaveReport(report *models.Report) error {
	if report.ID == 0 {
		if err := r.Create(report); err != nil {
			return fmt.Errorf("failed to create report: %w", err)
		}
	} else {
		if err := r.Update(report); err != nil {
			return fmt.Errorf("failed to update report: %w", err)
		}
	}
	return nil
}

func (r *ReportRepository) Create(report *models.Report) error {
	_, err := r.db.Exec(`
		INSERT INTO reports (order_id, engineer_id, has_repeat, repeat_date, repeat_note, description)
		VALUES (?, ?, ?, ?, ?, ?)`,
		report.OrderID, report.EngineerID, report.HasRepeat, report.RepeatDate, report.RepeatNote, report.Description,
	)
	return err
}

func (r *ReportRepository) Update(report *models.Report) error {
	_, err := r.db.Exec(`
		UPDATE reports SET has_repeat = ?, repeat_date = ?, repeat_note = ?, description = ?
		WHERE id = ?`,
		report.HasRepeat, report.RepeatDate, report.RepeatNote, report.Description, report.ID,
	)
	return err
}

func (r *ReportRepository) GetByOrderID(orderID string) (*models.Report, error) {
	row := r.db.QueryRow(`
		SELECT id, order_id, engineer_id, has_repeat, repeat_date, repeat_note, description FROM reports WHERE order_id = ?`, orderID)
	var report models.Report
	if err := row.Scan(&report.ID, &report.OrderID, &report.EngineerID, &report.HasRepeat, &report.RepeatDate, &report.RepeatNote, &report.Description); err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) GetCashReports(from, to string) ([]*domain.CashReport, error) {
	fromTime, _ := time.Parse("2006-01-02", from)
	toTime, _ := time.Parse("2006-01-02", to)
	toTime = toTime.AddDate(0, 0, 1)
	rows, err := r.db.Query(`
		SELECT
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
		FROM reports r
		LEFT JOIN orders o ON r.order_id = o.id
		JOIN engineers e ON r.engineer_id = e.id
		JOIN engineer_monthly_motivations mm
			  ON mm.engineer_id = o.engineer_id
			 AND DATE_FORMAT(o.scheduled_at,'%Y-%m') = DATE_FORMAT(mm.month,'%Y-%m')
		WHERE o.status in ('sent_to_cash', 'closed_finally') AND o.scheduled_at BETWEEN ? AND ?
		ORDER BY o.scheduled_at DESC

	`, fromTime, toTime)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*domain.CashReport

	for rows.Next() {
		var cr domain.CashReport
		err := rows.Scan(
			&cr.OrderID,
			&cr.ID,
			&cr.EngineerID,
			&cr.EngineerName,
			&cr.HasRepeat,
			&cr.RepeatDate,
			&cr.RepeatNote,
			&cr.Description,
			&cr.CreatedAt,
			&cr.ToCash,
			&cr.GaveCash,
			&cr.IssuedMoney,
			&cr.MotivationPercent,
			&cr.Prepayment,
			&cr.Salary,
			&cr.OrderStatus,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, &cr)
	}

	return reports, nil
}

func (r *ReportRepository) MarkIssued(orderId int64, gaveCash float64, issuedMoney float64, comment string) error {
	_, err := r.db.Exec(`
		UPDATE reports
		SET gave_cash = gave_cash + ? , issued_money = issued_money + ?, description = CONCAT_WS('\n', description, ?)
		WHERE order_id = ?
	`, gaveCash, issuedMoney, comment, orderId)

	return err
}

// CheckMotivationPrepaymentLimit report_repository.go
// Проверяем право приёма денег по мотивационному ограничению.
// Возвращает ошибку, если ограничение НЕ пройдено.
func (r *ReportRepository) CheckMotivationPrepaymentLimit(orderID int64) error {
	var info struct {
		MotivationPercent     float64
		TotalMotivationAmount float64
		PaidPrepayment        float64
	}

	err := r.db.QueryRow(`
	SELECT
		mm.motivation_percent,
		mm.total_motivation_amount,
		COALESCE(SUM(ep.paid_prepayment), 0) AS paid_prepayment
	FROM orders o
	JOIN engineer_monthly_motivations mm
		ON mm.engineer_id = o.engineer_id
		AND DATE_FORMAT(o.scheduled_at, '%Y-%m') = DATE_FORMAT(mm.month, '%Y-%m')
	LEFT JOIN engineer_payouts ep
		ON ep.engineer_id = o.engineer_id
	WHERE o.id = ?
	GROUP BY mm.motivation_percent, mm.total_motivation_amount
`, orderID).Scan(&info.MotivationPercent, &info.TotalMotivationAmount, &info.PaidPrepayment)

	if err != nil {
		return fmt.Errorf("failed to check motivation limit: %w", err)
	}

	// само правило
	if info.MotivationPercent >= 20 &&
		info.PaidPrepayment <= info.TotalMotivationAmount*0.5 {
		return fmt.Errorf("prepayment limit exceeded: paid %.2f, allowed %.2f (50%% of %.2f)",
			info.PaidPrepayment, info.TotalMotivationAmount*0.5, info.TotalMotivationAmount)
	}
	return nil
}
