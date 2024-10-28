package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FacilityRepository interface {
	GetAllFacilities() (*[]models.Facility, error)
	GetFacilityByID(id uuid.UUID) (*models.Facility, error)
}

type facilityRepository struct {
	db *gorm.DB
}

func NewFacilityRepository(db *gorm.DB) FacilityRepository {
	return &facilityRepository{db: db}
}

func (r *facilityRepository) GetAllFacilities() (*[]models.Facility, error) {
	var facilities []models.Facility
	if err := r.db.Find(&facilities).Error; err != nil {
		return nil, err
	}
	return &facilities, nil
}

func (r *facilityRepository) GetFacilityByID(id uuid.UUID) (*models.Facility, error) {
	var facility models.Facility
	if err := r.db.Where("id = ?", id).Find(&facility).Error; err != nil {
		return nil, err
	}
	return &facility, nil
}
