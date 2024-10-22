package repositories

import (
	"errors"
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomFacilityRepository interface {
	CreateRoomFacility(roomFacility *[]models.RoomFacility) error
	FindRoomFacilitiesByRoomID(id uuid.UUID) (*[]models.RoomFacility, error)
	UpdateRoomFacilityByRoomID(roomFacility *[]models.RoomFacility, id uuid.UUID) error
}

type roomFacilityRepository struct {
	db *gorm.DB
}

func NewRoomFacilityRepository(db *gorm.DB) RoomFacilityRepository {
	return &roomFacilityRepository{db: db}
}

func (r *roomFacilityRepository) CreateRoomFacility(roomFacility *[]models.RoomFacility) error {
	if err := r.db.Create(roomFacility).Error; err != nil {
		return err
	}
	return nil
}

func (r *roomFacilityRepository) FindRoomFacilitiesByRoomID(id uuid.UUID) (*[]models.RoomFacility, error) {
	var roomFacilities []models.RoomFacility
	if err := r.db.Where("rooming_id = ?", id).Find(&roomFacilities).Error; err != nil {
		return nil, err
	}
	return &roomFacilities, nil
}

func (r *roomFacilityRepository) UpdateRoomFacilityByRoomID(roomFacility *[]models.RoomFacility, id uuid.UUID) error {
	res := r.db.Delete(&roomFacility, "room_id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("room facility not found")
		}
	}

	if err := r.db.Create(roomFacility).Error; err != nil {
		return err
	}

	return nil
}
