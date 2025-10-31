package repositories

import (
	"erp/internal/app/models"
	"gorm.io/gorm"
)

type MotivationRepository struct {
	db *gorm.DB
}

func NewMotivationRepository(db *gorm.DB) *MotivationRepository {
	return &MotivationRepository{db: db}
}

func (r *MotivationRepository) Create(step *models.MotivationStep) error {
	return r.db.Create(step).Error
}

func (r *MotivationRepository) Update(step *models.MotivationStep) error {
	return r.db.Save(step).Error
}

func (r *MotivationRepository) Delete(id uint) error {
	return r.db.Delete(&models.MotivationStep{}, id).Error
}

func (r *MotivationRepository) GetAll() ([]models.MotivationStep, error) {
	var steps []models.MotivationStep
	err := r.db.Order("sort ASC").Find(&steps).Error
	return steps, err
}
