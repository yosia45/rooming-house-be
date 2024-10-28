package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeriodRepository interface {
	CreatePeriod(period *models.Period) error
	FindPeriodByName(name string) (*models.Period, error)
	FindPeriodByID(id uuid.UUID) (*models.Period, error)
	FindAllPeriods() (*[]models.Period, error)
	UpdatePeriodByID(period *models.Period, id uuid.UUID) error
	DeletePeriodByID(id uuid.UUID) error
}

type periodRepository struct {
	db *gorm.DB
}

func NewPeriodRepository(db *gorm.DB) PeriodRepository {
	return &periodRepository{db: db}
}

func (r *periodRepository) CreatePeriod(period *models.Period) error {
	if err := r.db.Create(period).Error; err != nil {
		return err
	}
	return nil
}

func (r *periodRepository) FindPeriodByName(name string) (*models.Period, error) {
	var period models.Period
	if err := r.db.Where("name LIKE ?", name).Find(&period).Error; err != nil {
		return nil, err
	}
	return &period, nil
}

func (r *periodRepository) FindPeriodByID(id uuid.UUID) (*models.Period, error) {
	var period models.Period
	if err := r.db.Where("id = ?", id).Find(&period).Error; err != nil {
		return nil, err
	}
	return &period, nil
}

func (r *periodRepository) FindAllPeriods() (*[]models.Period, error) {
	var periods []models.Period
	if err := r.db.Find(&periods).Error; err != nil {
		return nil, err
	}
	return &periods, nil
}

func (r *periodRepository) UpdatePeriodByID(period *models.Period, id uuid.UUID) error {
	res := r.db.Where("id = ?", id).Updates(period)
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

func (r *periodRepository) DeletePeriodByID(id uuid.UUID) error {
	res := r.db.Where("id = ?", id).Delete(&models.Period{})
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
