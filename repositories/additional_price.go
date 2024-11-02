package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdditionalPriceRepository interface {
	CreateAdditionalPrice(additionalPrice *models.AdditionalPrice) error
	FindAdditionalPriceByID(id uuid.UUID) (*models.AdditionalPriceResponse, error)
	FindAllAdditionalPrices(roomingHouseIDs []uuid.UUID) (*[]models.AdditionalPriceResponse, error)
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

func (r *additionalPriceRepository) FindAdditionalPriceByID(id uuid.UUID) (*models.AdditionalPriceResponse, error) {
	var additionalPrice models.AdditionalPrice
	if err := r.db.Preload("AdditionalPeriods").Preload("AdditionalPeriods.Period").Where("id = ?", id).First(&additionalPrice).Error; err != nil {
		return nil, err
	}

	response := models.AdditionalPriceResponse{
		ID:             additionalPrice.ID,
		Name:           additionalPrice.Name,
		RoomingHouseID: additionalPrice.RoomingHouseID,
		Prices:         make(map[string]float64),
	}

	for _, period := range additionalPrice.AdditionalPeriods {
		response.Prices[period.Period.Name] = period.Price
	}

	return &response, nil
}

func (r *additionalPriceRepository) FindAllAdditionalPrices(roomingHouseIDs []uuid.UUID) (*[]models.AdditionalPriceResponse, error) {
	var additionalPrices []models.AdditionalPrice

	for _, id := range roomingHouseIDs {
		var temp []models.AdditionalPrice
		if err := r.db.Preload("AdditionalPeriods").Preload("AdditionalPeriods.Period").Where("rooming_house_id = ?", id).Find(&temp).Error; err != nil {
			return nil, err
		}
		additionalPrices = append(additionalPrices, temp...)
	}

	var responses []models.AdditionalPriceResponse

	for _, price := range additionalPrices {
		response := models.AdditionalPriceResponse{
			ID:             price.ID,
			Name:           price.Name,
			RoomingHouseID: price.RoomingHouseID,
			Prices:         make(map[string]float64),
		}

		for _, period := range price.AdditionalPeriods {
			response.Prices[period.Period.Name] = period.Price
		}

		responses = append(responses, response)
	}

	return &responses, nil
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
