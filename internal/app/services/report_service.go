package services

import (
	"crypto/rand"
	"encoding/hex"
	"erp/internal/app/models"
	"erp/internal/app/repositories"
	"erp/internal/domain"
	"fmt"
	"gorm.io/gorm"
	"math"
	"os"
	"strconv"
	"time"
)

type ReportService struct {
	reportRepo           *repositories.ReportRepository
	orderRepo            *repositories.OrderRepository
	repeatRequestService *RepeatRequestService
	motivationCalculator *MotivationCalculator
}

type SubmitReportRequest struct {
	Token       string  `json:"token" binding:"required"`
	FinishPrice string  `json:"finish_price" binding:"required"`
	HasRepeat   bool    `json:"has_repeat"`
	RepeatDate  *string `json:"repeat_date,omitempty"` // формат "2006-01-02T15:04"
	RepeatNote  string  `json:"repeat_note,omitempty"`
	Description string  `json:"description,omitempty"`
}

func NewReportService(
	reportRepo *repositories.ReportRepository,
	orderRepo *repositories.OrderRepository,
	repeatRequestService *RepeatRequestService,
	db *gorm.DB,
) *ReportService {
	return &ReportService{
		reportRepo:           reportRepo,
		orderRepo:            orderRepo,
		repeatRequestService: repeatRequestService,
		motivationCalculator: &MotivationCalculator{DB: db},
	}
}

func (s *ReportService) GenerateReportLink(orderID, engineerID int64) (string, error) {
	b := make([]byte, 16)
	rand.Read(b)
	token := hex.EncodeToString(b)

	link := &models.ReportLink{
		Token:      token,
		OrderID:    orderID,
		EngineerID: engineerID,
		ExpiresAt:  time.Now().Add(24 * time.Hour), // действует сутки
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
	link, err := s.validateAndGetLink(req.Token)
	if err != nil {
		return err
	}

	report, err := s.createOrUpdateReport(link, req)
	if err != nil {
		return err
	}

	order, err := s.orderRepo.GetOrderByID(report.OrderID)
	if err != nil {
		return err
	}

	if !req.HasRepeat {
		order.Status = "closed_without_repeat"
	} else {
		//if err := s.handleRepeatLogic(req, report, order, hadRepeat); err != nil {
		//	return err
		//}
		if req.RepeatDate != nil {
			scheduledAt, err := time.ParseInLocation(
				"2006-01-02T15:04",
				*req.RepeatDate,
				time.Local,
			)
			if err != nil {
				return err
			}
			report.RepeatDate = &scheduledAt
		}
		err := s.repeatRequestService.Create(
			order.ID,
			uint(report.EngineerID),
			req.Description,
			*report.RepeatDate,
		)
		if err != nil {
			return err
		}
		totalOrderPrice, _ := strconv.ParseFloat(req.FinishPrice, 64)
		order.FinishPrice = fmt.Sprintf("%.2f", totalOrderPrice)
		order.AggregatorPayout = totalOrderPrice * (100 - order.OurPercent) / 100
		order.Status = "sent_to_cash"
	}

	// 5️⃣ Вычисление grossProfit
	totalOrderPrice, _ := strconv.ParseFloat(req.FinishPrice, 64)
	grossProfit := math.Round(totalOrderPrice * order.OurPercent / 100)
	order.FinishPrice = fmt.Sprintf("%.2f", totalOrderPrice)
	order.AggregatorPayout = totalOrderPrice - grossProfit

	isOrderRepeat := order.RepeatID != nil
	// 6️⃣ Расчет мотивации отчета
	if err := s.motivationCalculator.CalculateReportMotivation(report, totalOrderPrice, isOrderRepeat); err != nil {
		return err
	}

	if err := s.reportRepo.Update(report); err != nil {
		return err
	}

	// 7️⃣ Обновление месячного агрегата
	if err := s.motivationCalculator.RecalculateEngineerMonthlyMotivation(report); err != nil {
		return err
	}
	return s.orderRepo.UpdateOrder(order)
}

func (s *ReportService) validateAndGetLink(token string) (*models.ReportLink, error) {
	link, err := s.reportRepo.GetByToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if link.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	return link, nil
}

func (s *ReportService) createOrUpdateReport(link *models.ReportLink, req SubmitReportRequest) (*models.Report, error) {
	report, err := s.reportRepo.GetByOrderID(strconv.FormatInt(link.OrderID, 10))

	if err != nil || report == nil {
		// создаём новый
		report = &models.Report{
			OrderID:     link.OrderID,
			EngineerID:  link.EngineerID,
			HasRepeat:   req.HasRepeat,
			RepeatNote:  req.RepeatNote,
			Description: req.Description,
		}
	} else {
		// обновляем существующий
		report.HasRepeat = req.HasRepeat
		report.RepeatNote = req.RepeatNote
		report.Description = req.Description
	}

	if req.RepeatDate != nil {
		scheduledAt, err := time.ParseInLocation("2006-01-02T15:04", *req.RepeatDate, time.Local)
		if err != nil {
			return nil, err
		}
		report.RepeatDate = &scheduledAt
	}

	if err := s.reportRepo.SaveReport(report); err != nil {
		return nil, err
	}

	return report, nil
}

func (s *ReportService) handleRepeatLogic(req SubmitReportRequest, report *models.Report, order *models.Order, hadRepeat bool) error {
	// проверка был ли ранее создан повтор
	if hadRepeat {
		if err := s.updateRepeatOrder(report, order.ID); err != nil {
			return err
		}
	} else {
		//if err := s.createRepeatOrder(report, order); err != nil {
		//	return err
		//}
	}
	return nil
}

//
//func (s *ReportService) createRepeatOrder(report *models.Report, orig *models.Order) error {
//
//	nextErpNumber, err := s.orderRepo.GetNextERPNumber()
//	if err != nil {
//		return err
//	}
//
//	// Создаём новый заказ (повтор)
//	newOrder := *orig
//	//repeatID := orig.ID
//	//newOrder.RepeatID = &repeatID
//	newOrder.RepeatedBy = "engineer"
//	newOrder.RepeatDescription = report.RepeatNote
//	newOrder.ERPNumber = nextErpNumber
//	newOrder.Status = "working"
//
//	newOrder.ID = 0
//	newOrder.Aggregator = nil
//	newOrder.Problem = nil
//	newOrder.CreatedAt = time.Time{}
//	newOrder.UpdatedAt = time.Time{}
//
//	newOrder.ScheduledAt = *report.RepeatDate
//
//	return s.orderRepo.Create(&newOrder)
//}

func (s *ReportService) updateRepeatOrder(report *models.Report, RepeatID uint) error {

	repeatOrder, err := s.orderRepo.GetByRepeatID(RepeatID)
	if err != nil {
		return err
	}

	repeatOrder.RepeatDescription = report.RepeatNote
	repeatOrder.Note = fmt.Sprintf("Повтор от %s: %s", time.Now().Format("02.01.2006"), report.RepeatNote)
	repeatOrder.ScheduledAt = *report.RepeatDate

	return s.orderRepo.UpdateOrder(repeatOrder)
}

func (s *ReportService) GetCashReports(from, to string) ([]*domain.CashReport, error) {
	return s.reportRepo.GetCashReports(from, to)
}

func (s *ReportService) ReceiveCash(orderID int64, gaveCash float64, issuedMoney float64, comment string) error {
	if err := s.reportRepo.CheckMotivationPrepaymentLimit(orderID); err != nil {
		return err
	}

	if err := s.reportRepo.MarkIssued(orderID, gaveCash, issuedMoney, comment); err != nil {
		return err
	}
	return nil
}
