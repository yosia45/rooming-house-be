package repositories

import (
	"errors"
	"rooming-house-cms-be/models"

	"gorm.io/gorm"
)

type OwnerRepository interface {
	CreateOwner(user *models.Owner) error
	FindOwnerByEmail(email string) (*models.Owner, error)
}

type ownerRepository struct {
	db *gorm.DB
}

func NewOwnerRepository(db *gorm.DB) OwnerRepository {
	return &ownerRepository{db: db}
}

func (r *ownerRepository) CreateOwner(owner *models.Owner) error {
	if err := r.db.Create(owner).Error; err != nil {
		return err
	}
	return nil
}

func (r *ownerRepository) FindOwnerByEmail(email string) (*models.Owner, error) {
	var owner models.Owner
	if err := r.db.Where("email = ?", email).First(&owner).Error; err != nil {
		return nil, errors.New("owner not found")
	}
	return &owner, nil
}
