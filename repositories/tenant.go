package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantRepository interface {
	CreateTenant(tenant *models.Tenant) error
	FindAllTenants() (*[]models.Tenant, error)
	FindTenantByID(tenantID uuid.UUID) (*models.Tenant, error)
	UpdateTenantByID(tenant *models.Tenant, id uuid.UUID) error
	DeleteTenantByID(id uuid.UUID) error
}

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) CreateTenant(tenant *models.Tenant) error {
	if err := r.db.Create(tenant).Error; err != nil {
		return err
	}
	return nil
}

func (r *tenantRepository) FindAllTenants() (*[]models.Tenant, error) {
	var tenants []models.Tenant
	if err := r.db.Find(&tenants).Error; err != nil {
		return nil, err
	}
	return &tenants, nil
}

func (r *tenantRepository) FindTenantByID(tenantID uuid.UUID) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := r.db.Preload("AdditionalPrices", "tenant_additional_prices_id = ?", tenantID).Preload("AdditionalPrices.AdditionalPeriod").Where("id = ?", tenantID).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepository) UpdateTenantByID(tenant *models.Tenant, id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Updates(tenant).Error; err != nil {
		return err
	}
	return nil
}

func (r *tenantRepository) DeleteTenantByID(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&models.Tenant{}).Error; err != nil {
		return err
	}
	return nil
}
