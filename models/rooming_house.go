package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomingHouse struct {
	BaseModel
	Name                   string                 `json:"name" gorm:"not null"`
	Description            string                 `json:"description" gorm:"not null"`
	Address                string                 `json:"address" gorm:"not null"`
	FloorTotal             int                    `json:"floor_total" gorm:"not null"`
	OwnerID                uuid.UUID              `json:"owner_id" gorm:"not null; size:191"`
	Transactions           []Transaction          `json:"transactions" gorm:"foreignKey:RoomingHouseID"`
	RoomingHouseFacilities []RoomingHouseFacility `json:"rooming_house_facilities" gorm:"foreignKey:RoomingHouseID"`
	Rooms                  []Room                 `json:"rooms" gorm:"foreignKey:RoomingHouseID"`
	Admin                  Admin                  `json:"admin" gorm:"foreignKey:RoomingHouseID"`
}

type RoomingHouseBody struct {
	Name                    string      `json:"name" gorm:"not null"`
	Description             string      `json:"description" gorm:"not null"`
	Address                 string      `json:"address" gorm:"not null"`
	FloorTotal              int         `json:"floor_total" gorm:"not null"`
	UserID                  uuid.UUID   `json:"user_id" gorm:"not null"`
	RoomingHouseFacilityIDs []uuid.UUID `json:"rooming_house_facility_ids"`
}

func (rh *RoomingHouse) BeforeCreate(tx *gorm.DB) (err error) {
	rh.ID = uuid.New()
	rh.CreatedAt = time.Now()

	return
}
