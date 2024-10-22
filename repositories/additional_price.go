package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdditionalPriceRepository interface {
	CreateAdditionalPrice(additionalPrice *models.AdditionalPrice) error
	FindAdditionalPriceByID(id uuid.UUID) (*models.AdditionalPrice, error)
	FindAllAdditionalPrices() (*[]models.AdditionalPrice, error)
	UpdateAdditionalPriceByID(additionalPrice *models.AdditionalPrice, id uuid.UUID) error
	DeleteAdditionalPriceByID(id uuid.UUID) error
}

type additionalPriceRepository struct {
	db *gorm.DB
}

func NewAdditionalPriceRepository(db *gorm.DB) AdditionalPriceRepository {
	return &additionalPriceRepository{db: db}
}

func (r *additionalPriceRepository) CreateAdditionalPrice(additionalPrice *models.AdditionalPrice) error {
	if err := r.db.Create(additionalPrice).Error; err != nil {
		return err
	}
	return nil
}

func (r *additionalPriceRepository) FindAdditionalPriceByID(id uuid.UUID) (*models.AdditionalPrice, error) {
	var additionalPrice models.AdditionalPrice
	if err := r.db.Where("id = ?", id).First(&additionalPrice).Error; err != nil {
		return nil, err
	}
	return &additionalPrice, nil
}

func (r *additionalPriceRepository) FindAllAdditionalPrices() (*[]models.AdditionalPrice, error) {
	var additionalPrices []models.AdditionalPrice
	if err := r.db.Find(&additionalPrices).Error; err != nil {
		return nil, err
	}
	return &additionalPrices, nil
}

func (r *additionalPriceRepository) UpdateAdditionalPriceByID(additionalPrice *models.AdditionalPrice, id uuid.UUID) error {
	res := r.db.Where("id = ?", id).Updates(additionalPrice)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *additionalPriceRepository) DeleteAdditionalPriceByID(id uuid.UUID) error {
	res := r.db.Where("id = ?", id).Delete(&models.AdditionalPrice{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
