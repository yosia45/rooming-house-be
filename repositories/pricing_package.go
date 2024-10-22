package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PricingPackageRepository interface {
	CreatePricingPackage(pricingPackage *models.PricingPackage) error
	FindPricingPackageByID(id uuid.UUID) (*models.PricingPackage, error)
	FindAllPricingPackages() (*[]models.PricingPackage, error)
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

func (r *pricingPackageRepository) FindPricingPackageByID(id uuid.UUID) (*models.PricingPackage, error) {
	var pricingPackage models.PricingPackage
	if err := r.db.Where("id = ?", id).First(&pricingPackage).Error; err != nil {
		return nil, err
	}
	return &pricingPackage, nil
}

func (r *pricingPackageRepository) FindAllPricingPackages() (*[]models.PricingPackage, error) {
	var pricingPackages []models.PricingPackage
	if err := r.db.Find(&pricingPackages).Error; err != nil {
		return nil, err
	}
	return &pricingPackages, nil
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
