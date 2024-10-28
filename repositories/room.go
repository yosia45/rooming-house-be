package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepository interface {
	CreateRoom(room *models.Room) error
	// FindRoomsByRoomingHouseID(id uuid.UUID) (*[]models.Room, error)
	FindAllRooms(roomingHouseID uuid.UUID) (*[]models.Room, error)
	FindRoomByID(id uuid.UUID) (*models.Room, error)
	UpdateRoomByID(room *models.Room, id uuid.UUID) error
	DeleteRoomByID(id uuid.UUID) error
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) CreateRoom(room *models.Room) error {
	if err := r.db.Create(room).Error; err != nil {
		return err
	}
	return nil
}

// func (r *roomRepository) FindRoomsByRoomingHouseID(id uuid.UUID) (*[]models.Room, error) {
// 	var rooms []models.Room
// 	if err := r.db.Where("rooming_house_id = ?", id).Find(&rooms).Error; err != nil {
// 		return nil, err
// 	}
// 	return &rooms, nil
// }

func (r *roomRepository) FindAllRooms(roomingHouseID uuid.UUID) (*[]models.Room, error) {
	var rooms []models.Room
	if roomingHouseID == uuid.Nil {
		if err := r.db.Find(&rooms).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.db.Where("rooming_house_id = ?", roomingHouseID).Find(&rooms).Error; err != nil {
			return nil, err
		}
	}
	return &rooms, nil
}

func (r *roomRepository) FindRoomByID(id uuid.UUID) (*models.Room, error) {
	var room models.Room
	if err := r.db.Preload("Tenants").Preload("PricingPackages").Preload("Sizes").Preload("Facilities").Where("id = ?", id).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) UpdateRoomByID(room *models.Room, id uuid.UUID) error {
	res := r.db.Model(&models.Room{}).Where("id = ?", id).Updates(room)
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

func (r *roomRepository) DeleteRoomByID(id uuid.UUID) error {
	res := r.db.Delete(&models.Room{}, "id = ?", id)
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
