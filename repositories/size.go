package repositories

import (
	"errors"
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SizeRepository interface {
	CreateSize(size *models.Size) error
	FindSizeByID(id uuid.UUID) (*models.Size, error)
	FindAllSizes(roomingHouseID []uuid.UUID) (*[]models.AllSizeResponse, error)
	UpdateSizeByID(size *models.Size, id uuid.UUID) error
	DeleteSizeByID(id uuid.UUID) error
}

type sizeRepository struct {
	db *gorm.DB
}

func NewSizeRepository(db *gorm.DB) SizeRepository {
	return &sizeRepository{db: db}
}

func (r *sizeRepository) CreateSize(size *models.Size) error {
	if err := r.db.Create(size).Error; err != nil {
		return err
	}
	return nil
}

func (r *sizeRepository) FindSizeByID(id uuid.UUID) (*models.Size, error) {
	var size models.Size
	if err := r.db.Where("id = ?", id).First(&size).Error; err != nil {
		return nil, err
	}
	return &size, nil
}

func (r *sizeRepository) FindAllSizes(roomingHouseID []uuid.UUID) (*[]models.AllSizeResponse, error) {
	var sizes []models.Size

	if err := r.db.Where("rooming_house_id IN ?", roomingHouseID).Find(&sizes).Error; err != nil {
		return nil, err
	}

	var roomingHouses []models.RoomingHouse
	if err := r.db.Where("id IN ?", roomingHouseID).Find(&roomingHouses).Error; err != nil {
		return nil, err
	}

	roomingHouseMap := make(map[uuid.UUID]models.RoomingHouse)
	for _, rh := range roomingHouses {
		roomingHouseMap[rh.ID] = rh
	}

	var responses []models.AllSizeResponse
	for _, size := range sizes {
		response := models.AllSizeResponse{
			ID:    size.ID,
			Name:  size.Name,
			Width: size.Width,
			Long:  size.Long,
			RoomingHouse: models.TenantRoomingHouseResponse{
				ID:   roomingHouseMap[size.RoomingHouseID].ID,
				Name: roomingHouseMap[size.RoomingHouseID].Name,
			},
		}
		responses = append(responses, response)
	}

	return &responses, nil
}

func (r *sizeRepository) UpdateSizeByID(size *models.Size, id uuid.UUID) error {
	res := r.db.Where("id = ?", id).Updates(size)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("size not found")
		}
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("size not found")
	}

	return nil
}

func (r *sizeRepository) DeleteSizeByID(id uuid.UUID) error {
	res := r.db.Delete(&models.Size{}, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("size not found")
		}
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("size not found")
	}

	return nil
}
