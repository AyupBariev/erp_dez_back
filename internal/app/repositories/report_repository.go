package repositories

import (
	"database/sql"
	"erp/internal/app/models"
	"erp/internal/pkg/logger"
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
	_, err := r.db.Exec(`
		INSERT INTO reports (order_id, engineer_id, has_repeat, repeat_date, repeat_note, description)
		VALUES (?, ?, ?, ?, ?, ?)`,
		report.OrderID, report.EngineerID, report.HasRepeat, report.RepeatDate, report.RepeatNote, report.Description,
	)
	return err
}
