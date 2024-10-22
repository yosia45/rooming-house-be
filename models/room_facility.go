package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomFacility struct {
	BaseModel
	RoomID     uuid.UUID `json:"room_id" gorm:"not null;size:191"`
	FacilityID uuid.UUID `json:"facility_id" gorm:"not null;size:191"`
}

type RoomFacilityBody struct {
	RoomID     uuid.UUID `json:"room_id"`
	FacilityID uuid.UUID `json:"facility_id"`
}

func (rf *RoomFacility) BeforeCreate(tx *gorm.DB) (err error) {
	rf.ID = uuid.New()
	rf.CreatedAt = time.Now()

	return
}
