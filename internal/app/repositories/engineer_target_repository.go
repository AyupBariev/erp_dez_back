package repositories

import (
	"erp/internal/app/models"
	"gorm.io/gorm"
)

type EngineerTargetRepository struct {
	db *gorm.DB
}

func NewEngineerTargetRepository(db *gorm.DB) *EngineerTargetRepository {
	return &EngineerTargetRepository{db: db}
}

func (r *EngineerTargetRepository) Create(target *models.EngineerMotivationTarget) error {
	return r.db.Create(target).Error
}

func (r *EngineerTargetRepository) Update(target *models.EngineerMotivationTarget) error {
	return r.db.Save(target).Error
}

func (r *EngineerTargetRepository) GetByEngineerAndMonth(engineerID uint, monthStart, monthEnd string) ([]models.EngineerMotivationTarget, error) {
	var targets []models.EngineerMotivationTarget
	err := r.db.Where("engineer_id = ? AND target_month BETWEEN ? AND ?", engineerID, monthStart, monthEnd).Find(&targets).Error
	return targets, err
}
