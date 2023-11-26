package services

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"keyboard-analysis/internal/models"
)

type KeyboardService struct {
	db *gorm.DB
}

func NewKeyboardService(db *gorm.DB) *KeyboardService {
	return &KeyboardService{
		db: db,
	}
}

func (svc *KeyboardService) ProcessFlow(flows []models.KeyboardFlow) error {
	return svc.db.Save(flows).Error
}

func (svc *KeyboardService) History() ([]models.KeyboardFlow, error) {
	var res []models.KeyboardFlow
	if err := svc.query().Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *KeyboardService) SingleHistory(id uint) (*models.KeyboardFlow, error) {
	var res *models.KeyboardFlow
	if err := svc.query().Where("id = ?", id).First(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *KeyboardService) query() *gorm.DB {
	return svc.db.Preload(clause.Associations)
}
