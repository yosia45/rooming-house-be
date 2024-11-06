package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantRepository interface {
	CreateTenant(tenant *models.Tenant) error
	FindAllTenants(roomingHouseIDs []uuid.UUID, isTenant bool) (*[]models.AllTenantRepoResponse, error)
	FindTenantByID(tenantID uuid.UUID, roomingHouseIDs []uuid.UUID) (*models.TenantDetailResponse, error)
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

func (r *tenantRepository) FindAllTenants(roomingHouseIDs []uuid.UUID, isTenant bool) (*[]models.AllTenantRepoResponse, error) {

	var flatTenants []models.AllTenantRepoResponse

	query := r.db.Select("t.id, t.name, t.gender, t.start_date, t.end_date, r.id AS room_id, r.name AS room_name, rh.id AS rooming_house_id, rh.name AS rooming_house_name").
		Table("tenants t").
		Joins("JOIN rooms r ON t.room_id = r.id").
		Joins("JOIN rooming_houses rh ON t.rooming_house_id = rh.id").
		Where("t.rooming_house_id IN (?) AND t.deleted_at IS NULL", roomingHouseIDs)

	if isTenant {
		query = query.Where("t.is_tenant = ?", 1)
	}

	if err := query.Find(&flatTenants).Error; err != nil {
		return nil, err
	}

	return &flatTenants, nil
}

func (r *tenantRepository) FindTenantByID(tenantID uuid.UUID, roomingHouseIDs []uuid.UUID) (*models.TenantDetailResponse, error) {
	var tenantResponse models.TenantDetailResponse
	if err := r.db.
		Select("t.id, t.created_at, t.deleted_at, t.updated_at, t.name, t.gender, t.phone_number, t.emergency_contact, t.start_date, t.end_date, t.regular_payment_duration, t.is_tenant, t.is_deposit_paid, t.is_deposit_back, r.id AS room_id, r.name AS room_name, rh.id AS rooming_house_id, rh.name AS rooming_house_name, p.id AS period_id, p.name AS period_name").
		Table("tenants t").
		Joins("JOIN rooms r ON t.room_id = r.id").
		Joins("JOIN periods p ON t.period_id = p.id").
		Joins("JOIN rooming_houses rh ON rh.id = t.rooming_house_id").
		Where("t.id = ? AND t.deleted_at IS NULL AND t.rooming_house_id IN (?)", tenantID, roomingHouseIDs).
		First(&tenantResponse).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	var transactions []models.TransactionResponse
	if err := r.db.Table("transactions t").
		Select("t.id, t.created_at, t.day, t.month, t.year, t.amount, t.description, tc.name, tc.is_expense").
		Joins("JOIN transaction_categories tc ON t.transaction_category_id = tc.id").
		Where("tenant_id = ?", tenantID).
		Scan(&transactions).Error; err != nil {
		return nil, err
	}

	tenantResponse.Transactions = transactions

	var additionalPrices []models.AdditionalPriceDetail

	if err := r.db.Table("additional_prices ap").
		Select("ap.id, ap.name, adp.price, p.name AS period_name").
		Joins("JOIN tenant_additional_prices tap ON tap.additional_price_id = ap.id").
		Joins("JOIN additional_periods adp ON adp.additional_price_id = ap.id").
		Joins("JOIN periods p ON adp.period_id = p.id").
		Where("tap.tenant_id = ? AND p.id = ?", tenantID, tenantResponse.Period.ID).
		Scan(&additionalPrices).Error; err != nil {
		return nil, err
	}
	tenantResponse.AdditionalPrices = additionalPrices

	return &tenantResponse, nil
}

func (r *tenantRepository) UpdateTenantByID(tenant *models.Tenant, id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Updates(tenant).Error; err != nil {
		return err
	}
	return nil
}

func (r *tenantRepository) DeleteTenantByID(id uuid.UUID) error {
	res := r.db.Delete(&models.Tenant{}, "id = ?", id)
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
