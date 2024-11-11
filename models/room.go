package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	BaseModel
	Name           string        `json:"name" gorm:"not null"`
	Floor          int           `json:"floor" gorm:"not null"`
	MaxCapacity    int           `json:"max_capacity" gorm:"not null"`
	TenantID       uuid.UUID     `json:"tenant_id" gorm:"size:191"`
	SizeID         uuid.UUID     `json:"size_id" gorm:"not null;size:191"`
	PackageID      uuid.UUID     `json:"package_id" gorm:"not null;size:191"`
	RoomingHouseID uuid.UUID     `json:"rooming_house_id" gorm:"not null;size:191"`
	Facilities     []Facility    `gorm:"many2many:room_facilities;foreignKey:ID;joinForeignKey:RoomID;References:ID;joinReferences:FacilityID"`
	Tenants        *Tenant       `json:"tenant" gorm:"foreignKey:RoomID"`
	Transactions   []Transaction `json:"transactions" gorm:"foreignKey:RoomID"`
}

type AddRoomBody struct {
	Name           string      `json:"name"`
	Floor          int         `json:"floor"`
	MaxCapacity    int         `json:"max_capacity"`
	SizeID         uuid.UUID   `json:"size_id"`
	PackageID      uuid.UUID   `json:"package_id"`
	RoomingHouseID uuid.UUID   `json:"rooming_house_id"`
	RoomFacilities []uuid.UUID `json:"room_facilities"`
}

type UpdateRoomBody struct {
	Name           string      `json:"name"`
	Floor          int         `json:"floor"`
	MaxCapacity    int         `json:"max_capacity"`
	SizeID         uuid.UUID   `json:"size_id"`
	PackageID      uuid.UUID   `json:"package_id"`
	RoomFacilities []uuid.UUID `json:"room_facilities"`
}

type AllRoomResponse struct {
	ID             uuid.UUID            `json:"id" gorm:"column:room_id"`
	Name           string               `json:"name" gorm:"column:room_name"`
	Floor          int                  `json:"floor" gorm:"column:floor_number"`
	MaxCapacity    int                  `json:"max_capacity" gorm:"column:max_capacity"`
	RoomingHouseID uuid.UUID            `json:"rooming_house_id" gorm:"column:rooming_house_id"`
	Tenants        GetAllTenantResponse `json:"tenants"`
}

type RoomDetailResponse struct {
	ID             uuid.UUID                 `json:"id"`
	Name           string                    `json:"name"`
	Floor          int                       `json:"floor"`
	MaxCapacity    int                       `json:"max_capacity"`
	Size           Size                      `json:"size"`
	RoomingHouseID uuid.UUID                 `json:"rooming_house_id"`
	PricingPackage PackageResponse           `json:"pricing_package"`
	Tenants        *TenantRoomDetailResponse `json:"tenants"`
	Facilities     []Facility                `json:"facilities"`
}

type TenantRoomResponse struct {
	ID   uuid.UUID `json:"id" gorm:"column:room_id"`
	Name string    `json:"name" gorm:"column:room_name"`
}

func (r *Room) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	r.CreatedAt = time.Now()

	return
}
