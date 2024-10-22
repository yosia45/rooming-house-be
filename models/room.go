package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	BaseModel
	Name           string         `json:"name" gorm:"not null"`
	Floor          int            `json:"floor" gorm:"not null"`
	ExpiredDate    *time.Time     `json:"expired_date" gorm:"not null"`
	MaxCapacity    int            `json:"max_capacity" gorm:"not null"`
	StartDate      *time.Time     `json:"start_date" gorm:"not null"`
	IsVacant       bool           `json:"is_vacant" gorm:"not null"`
	SizeID         uuid.UUID      `json:"size_id" gorm:"not null;size:191"`
	PackageID      uuid.UUID      `json:"pricing_id" gorm:"not null;size:191"`
	RoomingHouseID uuid.UUID      `json:"rooming_house_id" gorm:"not null;size:191"`
	RoomFacility   []RoomFacility `json:"room_facility" gorm:"foreignKey:RoomID"`
	Tenants        []Tenant       `json:"tenant" gorm:"foreignKey:RoomID"`
	Transactions   []Transaction  `json:"transactions" gorm:"foreignKey:RoomID"`
}

type AddRoomBody struct {
	Name           string      `json:"name"`
	Floor          int         `json:"floor"`
	MaxCapacity    int         `json:"max_capacity"`
	SizeID         uuid.UUID   `json:"size_id"`
	PackageID      uuid.UUID   `json:"pricing_id"`
	RoomingHouseID uuid.UUID   `json:"rooming_house_id"`
	RoomFacilities []uuid.UUID `json:"room_facilities"`
}

type UpdateRoomBody struct {
	Name           string      `json:"name"`
	Floor          int         `json:"floor"`
	MaxCapacity    int         `json:"max_capacity"`
	SizeID         uuid.UUID   `json:"size_id"`
	PackageID      uuid.UUID   `json:"pricing_id"`
	RoomFacilities []uuid.UUID `json:"room_facilities"`
}

func (r *Room) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	r.CreatedAt = time.Now()

	return
}
