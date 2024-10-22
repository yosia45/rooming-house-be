package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomingHouseFacility struct {
	BaseModel
	RoomingHouseID uuid.UUID `json:"rooming_house_id" gorm:"not null;size:191"`
	FacilityID     uuid.UUID `json:"facility_id" gorm:"not null;size:191"`
}

type RoomingHouseFacilityBody struct {
	RoomingHouseID uuid.UUID `json:"rooming_house_id"`
	FacilityID     uuid.UUID `json:"facility_id"`
}

func (rhf *RoomingHouseFacility) BeforeCreate(tx *gorm.DB) (err error) {
	rhf.ID = uuid.New()
	rhf.CreatedAt = time.Now()

	return
}
