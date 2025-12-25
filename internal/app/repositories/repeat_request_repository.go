package repositories

import (
	"erp/internal/app/models"
	"errors"
	"gorm.io/gorm"
	"time"
)

type RepeatRequestRepository struct {
	db *gorm.DB
}

func NewRepeatRequestRepository(db *gorm.DB) *RepeatRequestRepository {
	return &RepeatRequestRepository{db: db}
}

// ========================
// CRUD
// ========================

func (r *RepeatRequestRepository) Create(req *models.RepeatRequest) error {
	return r.db.Create(req).Error
}

func (r *RepeatRequestRepository) GetByID(id uint) (*models.RepeatRequest, error) {
	var req models.RepeatRequest
	if err := r.db.
		Preload("Engineer").
		Preload("Order").
		First(&req, id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *RepeatRequestRepository) Update(req *models.RepeatRequest) error {
	return r.db.Save(req).Error
}

func (r *RepeatRequestRepository) List(status string) ([]models.RepeatRequest, error) {
	var reqs []models.RepeatRequest

	q := r.db.
		Preload("Order").
		Preload("Order.Problem").
		Preload("Order.Aggregator").
		Preload("Order.Engineer").
		Order("requested_at ASC")

	if status != "" {
		q = q.Where("status = ?", status)
	}

	if err := q.Find(&reqs).Error; err != nil {
		return nil, err
	}

	return reqs, nil
}

// ========================
// BUSINESS TRANSACTION
// ========================

func (r *RepeatRequestRepository) Confirm(
	requestID uint,
	curatorID int64,
	scheduledAt time.Time,
	repeatOrderID uint64,
) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		var req models.RepeatRequest
		if err := tx.First(&req, requestID).Error; err != nil {
			return err
		}

		if req.Status != "pending" {
			return errors.New("repeat request already processed")
		}

		now := time.Now()
		req.Status = "confirmed"
		req.Confirmed = true
		req.ConfirmedAt = &now
		req.ConfirmedBy = &curatorID
		req.RepeatOrderID = &repeatOrderID
		req.ScheduledAt = scheduledAt

		return tx.Save(&req).Error
	})
}
