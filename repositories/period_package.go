package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeriodPackageRepository interface {
	CreatePeriodPackage(periodPackage *[]models.PeriodPackage) error
	FindPeriodPackageByPackageID(packageID uuid.UUID) (*[]models.PeriodPackage, error)
	FindPeriodPackageByPeriodIDPackageID(periodID uuid.UUID, packageID uuid.UUID) (*models.PeriodPackage, error)
	UpdatePeriodPackageByPackageID(periodPackage []models.PeriodPackage, packageID uuid.UUID) error
}

type periodPackageRepository struct {
	db *gorm.DB
}

func NewPeriodPackageRepository(db *gorm.DB) PeriodPackageRepository {
	return &periodPackageRepository{db: db}
}

func (r *periodPackageRepository) CreatePeriodPackage(periodPackage *[]models.PeriodPackage) error {
	if err := r.db.Create(periodPackage).Error; err != nil {
		return err
	}
	return nil
}

func (r *periodPackageRepository) FindPeriodPackageByPackageID(packageID uuid.UUID) (*[]models.PeriodPackage, error) {
	var periodPackage []models.PeriodPackage
	if err := r.db.Find(&periodPackage).Where("pricing_package_id = ?", packageID).Error; err != nil {
		return nil, err
	}
	return &periodPackage, nil
}

func (r *periodPackageRepository) FindPeriodPackageByPeriodIDPackageID(periodID uuid.UUID, packageID uuid.UUID) (*models.PeriodPackage, error) {
	var periodPackage models.PeriodPackage
	if err := r.db.Where("period_id = ? AND pricing_package_id = ?", periodID, packageID).First(&periodPackage).Error; err != nil {
		return nil, err
	}
	return &periodPackage, nil
}

func (r *periodPackageRepository) UpdatePeriodPackageByPackageID(periodPackage []models.PeriodPackage, packageID uuid.UUID) error {
	for _, periodPackage := range periodPackage {
		err := r.db.Model(&models.PeriodPackage{}).
			Where("pricing_package_id = ? AND period_id = ?", packageID, periodPackage.PeriodID).
			Update("price", periodPackage.Price).Error

		if err != nil {
			return err
		}
	}

	return nil
}
