package repositories

import (
	"errors"
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateAdmin(user *models.Admin) error
	FindAdminByEmail(email string) (*models.Admin, error)
	FindAllAdmin(roomingHouseIDs []uuid.UUID) (*[]models.GetAllAdminResponse, error)
	DeleteAdminByID(id uuid.UUID) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) CreateAdmin(admin *models.Admin) error {
	if err := r.db.Create(admin).Error; err != nil {
		return err
	}
	return nil
}

func (r *adminRepository) FindAdminByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, errors.New("owner not found")
	}
	return &admin, nil
}

func (r *adminRepository) FindAllAdmin(roomingHouseIDs []uuid.UUID) (*[]models.GetAllAdminResponse, error) {
	var rawResults []struct {
		ID               uuid.UUID `json:"id"`
		FullName         string    `json:"full_name"`
		Username         string    `json:"username"`
		Role             string    `json:"role"`
		RoomingHouseID   uuid.UUID `json:"rooming_house_id"`
		RoomingHouseName string    `json:"rooming_house_name"`
	}

	var admins []models.GetAllAdminResponse

	if err := r.db.Table("admins").
		Select("admins.id, admins.full_name, admins.username, admins.role, admins.rooming_house_id AS rooming_house_id, rooming_houses.name AS rooming_house_name").
		Joins("JOIN rooming_houses ON admins.rooming_house_id = rooming_houses.id").
		Where("admins.rooming_house_id IN (?)", roomingHouseIDs).
		Scan(&rawResults).Error; err != nil {
		return nil, err
	}

	for _, rawResult := range rawResults {
		admins = append(admins, models.GetAllAdminResponse{
			ID:       rawResult.ID,
			FullName: rawResult.FullName,
			Username: rawResult.Username,
			Role:     rawResult.Role,
			RoomingHouse: models.TenantRoomingHouseResponse{
				ID:   rawResult.RoomingHouseID,
				Name: rawResult.RoomingHouseName,
			},
		})
	}

	return &admins, nil
}

func (r *adminRepository) DeleteAdminByID(id uuid.UUID) error {
	res := r.db.Where("id = ?", id).Delete(&models.Admin{})
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.New("admin not found")
		}
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("admin not found")
	}

	return nil
}
