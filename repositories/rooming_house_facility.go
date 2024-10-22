package repositories

import (
	"errors"
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomingHouseFacilityRepository interface {
	CreateRoomingHouseFacility(roomingHouseFacility *[]models.RoomingHouseFacility) error
	FindRoomingHouseFacilitiesByRoomingHouseID(id uuid.UUID) (*[]models.RoomingHouseFacility, error)
	UpdateRoomingHouseFacilityByRoomingHouseID(roomingHouseFacility *[]models.RoomingHouseFacility, id uuid.UUID) error
}

type roomingHouseFacilityRepository struct {
	db *gorm.DB
}

func NewRoomingHouseFacilityRepository(db *gorm.DB) RoomingHouseFacilityRepository {
	return &roomingHouseFacilityRepository{db: db}
}

func (r *roomingHouseFacilityRepository) CreateRoomingHouseFacility(roomingHouseFacility *[]models.RoomingHouseFacility) error {
	if err := r.db.Create(roomingHouseFacility).Error; err != nil {
		return err
	}
	return nil
}

func (r *roomingHouseFacilityRepository) FindRoomingHouseFacilitiesByRoomingHouseID(id uuid.UUID) (*[]models.RoomingHouseFacility, error) {
	var roomingHouseFacility []models.RoomingHouseFacility
	if err := r.db.Where("rooming_house_id = ?", id).Find(&roomingHouseFacility).Error; err != nil {
		return nil, err
	}
	return &roomingHouseFacility, nil
}

func (r *roomingHouseFacilityRepository) UpdateRoomingHouseFacilityByRoomingHouseID(roomingHouseFacility *[]models.RoomingHouseFacility, id uuid.UUID) error {
	res := r.db.Delete(&roomingHouseFacility, "rooming_house_id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("rooming house facility not found")
		}
	}

	if err := r.db.Create(roomingHouseFacility).Error; err != nil {
		return err
	}

	return nil
}
