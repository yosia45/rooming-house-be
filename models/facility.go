package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Facility struct {
	BaseModel
	Name                 string                 `json:"name" gorm:"not null"`
	Description          string                 `json:"description"`
	IsPublic             bool                   `json:"is_public" gorm:"not null"`
	IsRoom               bool                   `json:"is_room" gorm:"not null"`
	RoomFacility         []RoomFacility         `json:"room_facility" gorm:"foreignKey:FacilityID"`
	RoomingHouseFacility []RoomingHouseFacility `json:"rooming_house_facility" gorm:"foreignKey:FacilityID"`
}

type FacilityBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	IsRoom      bool   `json:"is_room"`
}

func (f *Facility) BeforeCreate(tx *gorm.DB) (err error) {
	f.ID = uuid.New()
	f.CreatedAt = time.Now()

	return
}
