package repositories

import (
	"errors"
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomingHouseRepository interface {
	CreateRoomingHouse(roomingHouse *models.RoomingHouse) error
	FindRoomingHouseByID(roomingHouseID uuid.UUID, userID uuid.UUID, role string) (*models.RoomingHouseByIDResponse, error)
	FindAllRoomingHouse(roomingHouseID uuid.UUID, userID uuid.UUID, role string) ([]models.AllRoomingHouseResponse, error)
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

func (r *roomingHouseRepository) FindRoomingHouseByID(roomingHouseID uuid.UUID, userID uuid.UUID, role string) (*models.RoomingHouseByIDResponse, error) {
	var roomingHouse models.RoomingHouse

	if err := r.db.Preload("Transactions").
		Preload("Facilities").
		Preload("Rooms").
		Preload("Admin").
		Where("id = ?", roomingHouseID).
		First(&roomingHouse).Error; err != nil {
		return nil, err
	}

	if role == "owner" {
		if roomingHouse.OwnerID != userID {
			return nil, errors.New("rooming house not found")
		}
	}

	roomingHouseResponse := models.RoomingHouseByIDResponse{
		ID:          roomingHouse.ID,
		Name:        roomingHouse.Name,
		Description: roomingHouse.Description,
		Address:     roomingHouse.Address,
		FloorTotal:  roomingHouse.FloorTotal,
		OwnerID:     roomingHouse.OwnerID,
		Admin: models.AdminResponse{
			ID:             roomingHouse.Admin.ID,
			Username:       roomingHouse.Admin.Username,
			Role:           roomingHouse.Admin.Role,
			RoomingHouseID: roomingHouse.Admin.RoomingHouseID,
		},
		Transactions: roomingHouse.Transactions,
		Facilities:   roomingHouse.Facilities,
		Rooms:        roomingHouse.Rooms,
	}

	return &roomingHouseResponse, nil
}

func (r *roomingHouseRepository) FindAllRoomingHouse(roomingHouseID uuid.UUID, userID uuid.UUID, role string) ([]models.AllRoomingHouseResponse, error) {
	var roomingHouses []models.RoomingHouse
	if role == "owner" {
		if err := r.db.Where("owner_id = ?", userID).Find(&roomingHouses).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.db.Where("id = ?", roomingHouseID).Find(&roomingHouses).Error; err != nil {
			return nil, err
		}
	}

	var roomingHouseResponses []models.AllRoomingHouseResponse
	for _, roomingHouse := range roomingHouses {
		roomingHouseResponse := models.AllRoomingHouseResponse{
			ID:          roomingHouse.ID,
			Name:        roomingHouse.Name,
			Description: roomingHouse.Description,
			Address:     roomingHouse.Address,
			FloorTotal:  roomingHouse.FloorTotal,
			OwnerID:     roomingHouse.OwnerID,
		}

		roomingHouseResponses = append(roomingHouseResponses, roomingHouseResponse)
	}

	return roomingHouseResponses, nil
}

func (r *roomingHouseRepository) Dashboard(roomingHouseIDs []uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.Where("rooming_house_id IN ?", roomingHouseIDs).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
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
