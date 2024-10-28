package repositories

import (
	"errors"
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdditionalPeriodRepository interface {
	CreateAdditionalPeriod(additionalPeriod *[]models.AdditionalPeriod) error
	FindPrice(periodID uuid.UUID, additionalPriceID uuid.UUID) (*models.AdditionalPeriod, error)
	UpdateAdditionalPeriod(additionalPeriod *[]models.AdditionalPeriod, additionalPriceID uuid.UUID) error
}

type additionalPeriodRepository struct {
	db *gorm.DB
}

func NewAdditionalPeriodRepository(db *gorm.DB) AdditionalPeriodRepository {
	return &additionalPeriodRepository{db: db}
}

func (r *additionalPeriodRepository) CreateAdditionalPeriod(additionalPeriod *[]models.AdditionalPeriod) error {
	if err := r.db.Create(additionalPeriod).Error; err != nil {
		return err
	}
	return nil
}

func (r *additionalPeriodRepository) FindPrice(periodID uuid.UUID, additionalPriceID uuid.UUID) (*models.AdditionalPeriod, error) {
	var additionalPeriod models.AdditionalPeriod
	if err := r.db.Where("period_id = ? AND additional_price_id = ?", periodID, additionalPriceID).First(&additionalPeriod).Error; err != nil {
		return nil, err
	}
	return &additionalPeriod, nil
}

func (r *additionalPeriodRepository) UpdateAdditionalPeriod(additionalPeriod *[]models.AdditionalPeriod, additionalPriceID uuid.UUID) error {
	res := r.db.Delete(&additionalPeriod, "additional_price_id = ?", additionalPriceID)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("additional period not found")
		}
	}

	if err := r.db.Create(additionalPeriod).Error; err != nil {
		return err
	}

	return nil
}
