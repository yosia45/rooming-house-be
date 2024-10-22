package repositories

import (
	"errors"
	"rooming-house-cms-be/models"

	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateAdmin(user *models.Admin) error
	FindAdminByEmail(email string) (*models.Admin, error)
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
