package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantAdditionalRepository interface {
	CreateTenantAdditional(tenantAdditional *[]models.TenantAdditionalPrice) error
	FindAllTenantAdditionalsByTenantID(id uuid.UUID) (*[]models.TenantAdditionalPrice, error)
	UpdateTenantAdditionalByTenantID(tenantAdditional *[]models.TenantAdditionalPrice, id uuid.UUID) error
}

type tenantAdditionalRepository struct {
	db *gorm.DB
}

func NewTenantAdditionalRepository(db *gorm.DB) TenantAdditionalRepository {
	return &tenantAdditionalRepository{db: db}
}

func (r *tenantAdditionalRepository) CreateTenantAdditional(tenantAdditional *[]models.TenantAdditionalPrice) error {
	if err := r.db.Create(tenantAdditional).Error; err != nil {
		return err
	}
	return nil
}

func (r *tenantAdditionalRepository) FindAllTenantAdditionalsByTenantID(id uuid.UUID) (*[]models.TenantAdditionalPrice, error) {
	var tenantAdditionals []models.TenantAdditionalPrice
	if err := r.db.Where("tenant_id = ?", id).Find(&tenantAdditionals).Error; err != nil {
		return nil, err
	}

	return &tenantAdditionals, nil
}

func (r *tenantAdditionalRepository) UpdateTenantAdditionalByTenantID(tenantAdditional *[]models.TenantAdditionalPrice, id uuid.UUID) error {
	res := r.db.Delete(&tenantAdditional, "tenant_id = ?", id)
	if res.Error != nil {
		return res.Error
	}

	if err := r.db.Create(tenantAdditional).Error; err != nil {
		return err
	}

	return nil
}
