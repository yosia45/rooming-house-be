package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PricingPackageRepository interface {
	CreatePricingPackage(pricingPackage *models.PricingPackage) error
	FindPricingPackageByID(packageID uuid.UUID) (*models.PricingPackage, error)
	FindAllPricingPackages(roomingHouseIDs []uuid.UUID) (*[]models.PackageResponse, error)
	UpdatePricingPackageByID(pricingPackage *models.PricingPackage, id uuid.UUID) error
	DeletePricingPackageByID(id uuid.UUID) error
}

type pricingPackageRepository struct {
	db *gorm.DB
}

func NewPricingPackageRepository(db *gorm.DB) PricingPackageRepository {
	return &pricingPackageRepository{db: db}
}

func (r *pricingPackageRepository) CreatePricingPackage(pricingPackage *models.PricingPackage) error {
	if err := r.db.Create(pricingPackage).Error; err != nil {
		return err
	}
	return nil
}

func (r *pricingPackageRepository) FindPricingPackageByID(packageID uuid.UUID) (*models.PricingPackage, error) {
	var pricingPackage models.PricingPackage

	if err := r.db.Where("id = ?", packageID).First(&pricingPackage).Error; err != nil {
		return nil, err
	}
	return &pricingPackage, nil
}

func (r *pricingPackageRepository) FindAllPricingPackages(roomingHouseIDs []uuid.UUID) (*[]models.PackageResponse, error) {
	var pricingPackages []models.PricingPackage

	for _, id := range roomingHouseIDs {
		var temp []models.PricingPackage
		if err := r.db.Preload("PeriodPackages").Preload("PeriodPackages.Period").Where("rooming_house_id = ?", id).Find(&temp).Error; err != nil {
			return nil, err
		}
		pricingPackages = append(pricingPackages, temp...)
	}

	var responses []models.PackageResponse

	for _, pkg := range pricingPackages {
		response := models.PackageResponse{
			ID:             pkg.ID,
			Name:           pkg.Name,
			RoomingHouseID: pkg.RoomingHouseID,
			Prices:         make(map[string]float64),
		}

		// Memetakan harga berdasarkan unit periode
		for _, periodPackage := range pkg.PeriodPackages {
			response.Prices[periodPackage.Period.Unit] = periodPackage.Price
		}

		responses = append(responses, response)
	}

	return &responses, nil
}

func (r *pricingPackageRepository) UpdatePricingPackageByID(pricingPackage *models.PricingPackage, id uuid.UUID) error {
	res := r.db.Where("id = ?", id).Updates(pricingPackage)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return res.Error
		}
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *pricingPackageRepository) DeletePricingPackageByID(id uuid.UUID) error {
	res := r.db.Where("id = ?", id).Delete(&models.PricingPackage{})
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return res.Error
		}
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
