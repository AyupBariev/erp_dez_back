package services

import (
	"crypto/rand"
	"encoding/hex"
	"erp/internal/app/models"
	"erp/internal/app/repositories"
	"fmt"
	"gorm.io/gorm"
	"os"
	"time"
)

type ReportService struct {
	reportRepo           *repositories.ReportRepository
	orderRepo            *repositories.OrderRepository
	motivationCalculator *MotivationCalculator
}

type SubmitReportRequest struct {
	Token       string  `json:"token" binding:"required"`
	HasRepeat   bool    `json:"has_repeat"`
	RepeatDate  *string `json:"repeat_date,omitempty"` // формат "2006-01-02T15:04"
	RepeatNote  string  `json:"repeat_note,omitempty"`
	Description string  `json:"description,omitempty"`
}

func NewReportService(reportRepo *repositories.ReportRepository, orderRepo *repositories.OrderRepository, db *gorm.DB) *ReportService {
	return &ReportService{reportRepo: reportRepo, orderRepo: orderRepo, motivationCalculator: &MotivationCalculator{DB: db}}

}

func (s *ReportService) GenerateReportLink(orderID, engineerID int64) (string, error) {
	b := make([]byte, 16)
	rand.Read(b)
	token := hex.EncodeToString(b)

	link := &models.ReportLink{
		Token:      token,
		OrderID:    orderID,
		EngineerID: engineerID,
		ExpiresAt:  time.Now().Add(1 * time.Hour),
	}

	if err := s.reportRepo.SaveReportLink(link); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/reports/submit?token=%s", os.Getenv("ERP_FRONT_URI"), token), nil
}

func (s *ReportService) GetByToken(token string) (*models.ReportLink, error) {
	return s.reportRepo.GetByToken(token)
}

func (s *ReportService) SubmitReport(req SubmitReportRequest) error {
	link, err := s.reportRepo.GetByToken(req.Token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	if link.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("token expired")
	}

	report := &models.Report{
		OrderID:     link.OrderID,
		EngineerID:  link.EngineerID,
		HasRepeat:   req.HasRepeat,
		RepeatNote:  req.RepeatNote,
		Description: req.Description,
	}

	if req.RepeatDate != nil {
		if t, err := time.Parse("2006-01-02T15:04", *req.RepeatDate); err == nil {
			report.RepeatDate = &t
		}
	}

	if err := s.reportRepo.SaveReport(report); err != nil {
		return err
	}

	// Получаем исходный заказ
	order, err := s.orderRepo.GetOrderByErpNumber(report.OrderID)
	if err != nil {
		return err
	}
	isRepeat := req.HasRepeat
	if err := s.motivationCalculator.UpdateEngineerMonthlyMotivation(
		report.EngineerID,
		order.Price,
		order.OurPercent,
		isRepeat,
	); err != nil {
		return err
	}

	// Логика статусов и повтора:
	if isRepeat {
		// 1️⃣ Создаём повтор
		if err := s.createRepeatOrder(report, order); err != nil {
			return err
		}
		// 2️⃣ Закрываем исходный заказ
		return s.orderRepo.UpdateStatus(report.OrderID, "closed_finally")
	}

	// Без повтора — просто закрываем
	return s.orderRepo.UpdateStatus(report.OrderID, "closed_without_repeat")
}

func (s *ReportService) createRepeatOrder(report *models.Report, orig *models.Order) error {

	nextErpNumber, err := s.orderRepo.GetNextERPNumber()
	if err != nil {
		return err
	}

	// Создаём новый заказ (повтор)
	newOrder := *orig
	newOrder.RepeatID = orig.ID
	newOrder.RepeatedBy = "engineer"
	newOrder.RepeatDescription = report.RepeatNote
	newOrder.ERPNumber = nextErpNumber
	newOrder.Status = "new"
	newOrder.Note = fmt.Sprintf("Повтор от %s: %s", time.Now().Format("02.01.2006"), report.RepeatNote)

	if report.RepeatDate != nil {
		newOrder.ScheduledAt = *report.RepeatDate
	}

	return s.orderRepo.Create(&newOrder)
}
