package repositories

import (
	"errors"
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomingHouseRepository interface {
	CreateRoomingHouse(roomingHouse *models.RoomingHouse) error
	FindRoomingHouseByID(id uuid.UUID) (*models.RoomingHouse, error)
	FindAllRoomingHouse(roomingHouseID uuid.UUID, userID uuid.UUID, role string) ([]models.RoomingHouse, error)
	UpdateRoomingHouse(roomingHouse *models.RoomingHouse, id uuid.UUID) error
	DeleteRoomingHouse(id uuid.UUID) error
}

type roomingHouseRepository struct {
	db *gorm.DB
}

func NewRoomingHouseRepository(db *gorm.DB) RoomingHouseRepository {
	return &roomingHouseRepository{db: db}
}

func (r *roomingHouseRepository) CreateRoomingHouse(roomingHouse *models.RoomingHouse) error {
	if err := r.db.Create(roomingHouse).Error; err != nil {
		return err
	}
	return nil
}

func (r *roomingHouseRepository) FindRoomingHouseByID(id uuid.UUID) (*models.RoomingHouse, error) {
	var roomingHouse models.RoomingHouse
	if err := r.db.Preload("Transactions").
		Preload("RoomingHouseFacilities").
		Preload("Rooms").
		Preload("Admin").
		Where("id = ?", id).
		First(&roomingHouse).Error; err != nil {
		return nil, err
	}
	return &roomingHouse, nil
}

func (r *roomingHouseRepository) FindAllRoomingHouse(roomingHouseID uuid.UUID, userID uuid.UUID, role string) ([]models.RoomingHouse, error) {
	var roomingHouses []models.RoomingHouse
	if role == "owner" {
		if err := r.db.Find(&roomingHouses).Where("owner_id = ?", userID).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.db.Find(&roomingHouses).Where("id = ?", roomingHouseID).Error; err != nil {
			return nil, err
		}
	}
	return roomingHouses, nil
}

func (r *roomingHouseRepository) UpdateRoomingHouse(roomingHouse *models.RoomingHouse, id uuid.UUID) error {
	res := r.db.Model(&roomingHouse).Where("id = ?", id).Updates(roomingHouse)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("rooming house not found")
		}

		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("rooming house not found")
	}

	return nil
}

func (r *roomingHouseRepository) DeleteRoomingHouse(id uuid.UUID) error {
	roomingHouse := models.RoomingHouse{}

	res := r.db.Delete(&roomingHouse, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("rooming house not found")
		}

		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("rooming house not found")
	}

	return nil
}
