package repositories

import (
	"fmt"
	"rooming-house-cms-be/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantRepository interface {
	CreateTenant(tenant *models.Tenant) error
	FindAllTenants(roomingHouseIDs []uuid.UUID, IsTenant bool) (*[]models.AllTenantRepoResponse, error)
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

func (r *tenantRepository) FindAllTenants(roomingHouseIDs []uuid.UUID, IsTenant bool) (*[]models.AllTenantRepoResponse, error) {
	var flatTenants []models.AllTenantRepoResponse

	query := r.db.Select("t.id, t.name, t.gender, t.start_date, t.end_date, t.is_tenant, r.id AS room_id, r.name AS room_name, rh.id AS rooming_house_id, rh.name AS rooming_house_name").
		Table("tenants t").
		Joins("LEFT JOIN rooms r ON t.room_id = r.id").
		Joins("JOIN rooming_houses rh ON t.rooming_house_id = rh.id").
		Where("t.rooming_house_id IN (?) AND t.deleted_at IS NULL", roomingHouseIDs)

	if IsTenant {
		query = query.Where("t.is_tenant = true")
	}

	if err := query.Find(&flatTenants).Error; err != nil {
		return nil, err
	}

	return &flatTenants, nil
}

func (r *tenantRepository) FindTenantByID(tenantID uuid.UUID, roomingHouseIDs []uuid.UUID) (*models.TenantDetailResponse, error) {
	now := time.Now()

	var tenant models.Tenant
	// Query tenant utama untuk mengecek nilai is_tenant
	if err := r.db.
		Table("tenants t").
		Select("t.id, t.name, t.gender, t.phone_number, t.emergency_contact, t.is_tenant, t.rooming_house_id").
		Where("t.id = ? AND t.deleted_at IS NULL AND t.rooming_house_id IN (?)", tenantID, roomingHouseIDs).
		First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tenant with ID %s not found", tenantID)
		}
		return nil, err
	}

	// Jika is_tenant = false, isi hanya data terbatas
	if !tenant.IsTenant {
		var roomingHouse models.TenantRoomingHouseResponse
		// Ambil informasi Rooming House
		if err := r.db.
			Table("rooming_houses").
			Select("id, name").
			Where("id = ?", tenant.RoomingHouseID).
			Scan(&roomingHouse).Error; err != nil {
			return nil, err
		}

		// Buat response sederhana
		return &models.TenantDetailResponse{
			ID:               tenant.ID,
			Name:             tenant.Name,
			Gender:           tenant.Gender,
			PhoneNumber:      tenant.PhoneNumber,
			EmergencyContact: tenant.EmergencyContact,
			RoomingHouse:     roomingHouse,
		}, nil
	}

	// Jika is_tenant = true, ambil seluruh data detail tenant
	var tenantResponse models.TenantDetailResponse
	if err := r.db.
		Select("t.id, t.created_at, t.deleted_at, t.updated_at, t.name, t.gender, t.phone_number, t.emergency_contact, t.room_id as booked_room_id, t.start_date, t.end_date, t.regular_payment_duration, t.is_tenant, t.is_deposit_paid, t.is_deposit_back, r.id AS room_id, r.name AS room_name, rh.id AS rooming_house_id, rh.name AS rooming_house_name, p.id AS period_id, p.name AS period_name").
		Table("tenants t").
		Joins("LEFT JOIN rooms r ON t.room_id = r.id AND t.start_date <= ? AND t.end_date >= ?", now, now).
		Joins("JOIN periods p ON t.period_id = p.id").
		Joins("JOIN rooming_houses rh ON rh.id = t.rooming_house_id").
		Where("t.id = ? AND t.deleted_at IS NULL AND t.rooming_house_id IN (?)", tenantID, roomingHouseIDs).
		First(&tenantResponse).Error; err != nil {
		return nil, err
	}

	// Tenant Assists
	var tenantAssists []models.TenantAssistResponse
	if err := r.db.
		Table("tenants t").
		Select("ta.id, ta.name, ta.gender, ta.phone_number, ta.is_tenant, ta.tenant_id, ta.rooming_house_id").
		Joins("JOIN tenants ta ON t.id = ta.tenant_id").
		Where("t.id = ? AND ta.is_tenant = false AND ta.deleted_at IS NULL", tenantID).
		Scan(&tenantAssists).Error; err != nil {
		return nil, err
	}
	tenantResponse.TenantAssists = tenantAssists

	// Transactions
	var transactions []models.TransactionResponse
	if err := r.db.Table("transactions t").
		Select("t.id, t.created_at, t.day, t.month, t.year, t.amount, t.description, tc.name, tc.is_expense").
		Joins("JOIN transaction_categories tc ON t.transaction_category_id = tc.id").
		Where("tenant_id = ?", tenantID).
		Scan(&transactions).Error; err != nil {
		return nil, err
	}
	tenantResponse.Transactions = transactions

	// Additional Prices
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
