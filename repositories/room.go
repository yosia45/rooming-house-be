package repositories

import (
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepository interface {
	CreateRoom(room *models.Room) error
	FindAllRooms(roomingHouseIDs []uuid.UUID) (*[]models.AllRoomResponse, error)
	FindRoomByID(roomID uuid.UUID, roomingHouseID uuid.UUID, userID uuid.UUID, userRole string) (*models.RoomDetailResponse, error)
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

func (r *roomRepository) FindAllRooms(roomingHouseIDs []uuid.UUID) (*[]models.AllRoomResponse, error) {
	var response []models.AllRoomResponse

	now := time.Now()

	query := `
		SELECT 
			r.id AS room_id,
			r.name AS room_name,
			r.floor AS floor_number,
			r.max_capacity,
			r.rooming_house_id,
			t.id AS tenant_id,
			t.name AS tenant_name,
			t.gender AS tenant_gender,
			t.start_date AS tenant_start_date,
			t.end_date AS tenant_end_date,
			t.room_id AS tenant_room_id,
			t.rooming_house_id AS tenant_rooming_house_id
		FROM 
			rooms r
		LEFT JOIN 
			tenants t
		ON 
			r.id = t.room_id 
			AND t.is_tenant = 1 
			AND t.start_date <= ? 
			AND t.end_date >= ?
		WHERE 
			r.rooming_house_id IN (?);
	`

	// Replace with raw query and scan into a struct
	rows, err := r.db.Raw(query, now, now, roomingHouseIDs).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over rows and map to response
	for rows.Next() {
		var room models.AllRoomResponse
		var tenant struct {
			ID             *uuid.UUID `json:"id"`
			Name           *string    `json:"name"`
			Gender         *string    `json:"gender"`
			StartDate      *time.Time `json:"start_date"`
			EndDate        *time.Time `json:"end_date"`
			RoomID         *uuid.UUID `json:"room_id"`
			RoomingHouseID *uuid.UUID `json:"rooming_house_id"`
		}

		// Scan each column to the correct struct
		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.Floor,
			&room.MaxCapacity,
			&room.RoomingHouseID,
			&tenant.ID,
			&tenant.Name,
			&tenant.Gender,
			&tenant.StartDate,
			&tenant.EndDate,
			&tenant.RoomID,
			&tenant.RoomingHouseID,
		)
		if err != nil {
			return nil, err
		}

		// Check if tenant ID is NULL
		if tenant.ID == nil {
			room.Tenants = models.GetAllTenantResponse{} // Empty tenant object
		} else {
			room.Tenants = models.GetAllTenantResponse{
				ID:             *tenant.ID,
				Name:           utils.PtrToString(tenant.Name),
				Gender:         utils.PtrToString(tenant.Gender),
				StartDate:      tenant.StartDate,
				EndDate:        tenant.EndDate,
				RoomID:         *tenant.RoomID,
				RoomingHouseID: *tenant.RoomingHouseID,
			}
		}

		response = append(response, room)
	}

	return &response, nil
}

func (r *roomRepository) FindRoomByID(roomID uuid.UUID, roomingHouseID uuid.UUID, userPayload uuid.UUID, userRole string) (*models.RoomDetailResponse, error) {
	var room models.Room
	now := time.Now()

	if userRole == "admin" {
		if err := r.db.Preload("Facilities").Where("id = ? AND rooming_house_id = ?", roomID, roomingHouseID).First(&room).Error; err != nil {
			return nil, err
		}
	} else {
		var roomingHouses []models.RoomingHouse

		if err := r.db.Find(&roomingHouses, "owner_id = ?", userPayload).Error; err != nil {
			return nil, err
		}

		roomingHouseIDs := make([]uuid.UUID, len(roomingHouses))
		for i, rh := range roomingHouses {
			roomingHouseIDs[i] = rh.ID
		}

		if err := r.db.Preload("Facilities").Where("id = ? AND rooming_house_id IN ?", roomID, roomingHouseIDs).First(&room).Error; err != nil {
			return nil, err
		}
	}

	var tenantWithAssists models.TenantRoomDetailResponse

	// Query pertama untuk tenant utama
	if err := r.db.Table("tenants as t").
		Select("t.id, t.name, t.gender, t.phone_number, t.emergency_contact, t.is_tenant, t.start_date, t.end_date, t.regular_payment_duration").
		Where("t.room_id = ? AND t.is_tenant = true AND t.start_date <= ? AND t.end_date >= ? AND t.deleted_at IS NULL", roomID, now, now).
		Scan(&tenantWithAssists).Error; err != nil {
		return nil, err
	}

	// Query kedua untuk assist
	var tenantAssists []models.TenantAssistResponse
	if err := r.db.Table("tenants as ta").
		Select("ta.id as id, ta.name as name, ta.gender as gender, ta.phone_number as phoneNumber, ta.is_tenant as isTenant").
		Where("ta.tenant_id = ? AND ta.is_tenant = false AND ta.deleted_at IS NULL", tenantWithAssists.ID).
		Scan(&tenantAssists).Error; err != nil {
		return nil, err
	}

	// Gabungkan hasil
	tenantWithAssists.TenantAssists = tenantAssists

	var size models.Size
	if err := r.db.Where("id = ?", room.SizeID).First(&size).Error; err != nil {
		return nil, err
	}

	var pricingPackage models.PricingPackage
	if err := r.db.Preload("PeriodPackages.Period").
		Where("id = ?", room.PackageID).
		First(&pricingPackage).Error; err != nil {
		return nil, err
	}

	roomDetailResponse := models.RoomDetailResponse{
		ID:             room.ID,
		Name:           room.Name,
		Floor:          room.Floor,
		MaxCapacity:    room.MaxCapacity,
		Size:           size,
		RoomingHouseID: room.RoomingHouseID,
		Tenants:        &tenantWithAssists,
		Facilities:     room.Facilities,
		PricingPackage: models.PackageResponse{
			ID:             pricingPackage.ID,
			Name:           pricingPackage.Name,
			RoomingHouseID: pricingPackage.RoomingHouseID,
			Prices:         map[string]float64{},
		},
	}

	for _, periodPackage := range pricingPackage.PeriodPackages {
		roomDetailResponse.PricingPackage.Prices[periodPackage.Period.Name] = periodPackage.Price
	}

	return &roomDetailResponse, nil
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
