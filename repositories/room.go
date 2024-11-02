package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepository interface {
	CreateRoom(room *models.Room) error
	// FindRoomsByRoomingHouseID(id uuid.UUID) (*[]models.Room, error)
	FindAllRooms(roomingHouseID uuid.UUID) (*[]models.AllRoomResponse, error)
	FindRoomByID(id uuid.UUID) (*models.RoomDetailResponse, error)
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

func (r *roomRepository) FindAllRooms(roomingHouseID uuid.UUID) (*[]models.AllRoomResponse, error) {
	var rooms []models.Room
	if roomingHouseID == uuid.Nil {
		if err := r.db.Preload("Tenants").Find(&rooms).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.db.Preload("Tenants").Where("rooming_house_id = ?", roomingHouseID).Find(&rooms).Error; err != nil {
			return nil, err
		}
	}

	// Mempersiapkan slice untuk response
	var response []models.AllRoomResponse

	// Memetakan room ke response
	for _, room := range rooms {
		tenantResponses := make([]models.GetAllTenantResponse, len(room.Tenants))
		for i, tenant := range room.Tenants {
			tenantResponses[i] = models.GetAllTenantResponse{
				ID:             tenant.ID,
				Name:           tenant.Name,
				Gender:         tenant.Gender,
				StartDate:      tenant.StartDate,
				EndDate:        tenant.EndDate,
				RoomID:         tenant.RoomID,
				RoomingHouseID: tenant.RoomingHouseID,
			}
		}

		response = append(response, models.AllRoomResponse{
			ID:             room.ID,
			Name:           room.Name,
			Floor:          room.Floor,
			MaxCapacity:    room.MaxCapacity,
			IsVacant:       room.IsVacant,
			Tenants:        tenantResponses,
			RoomingHouseID: room.RoomingHouseID,
		})
	}

	return &response, nil
}

func (r *roomRepository) FindRoomByID(id uuid.UUID) (*models.RoomDetailResponse, error) {
	var room models.Room
	if err := r.db.Preload("Tenants").Preload("Facilities").Where("id = ?", id).First(&room).Error; err != nil {
		return nil, err
	}

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
		IsVacant:       room.IsVacant,
		Size:           size,
		RoomingHouseID: room.RoomingHouseID,
		Tenants:        []models.Tenant{},
		Facilities:     []models.Facility{},
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
