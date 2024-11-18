package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Size struct {
	BaseModel
	Name           string    `json:"name" gorm:"not null"`
	Width          float64   `json:"width" gorm:"not null"`
	Long           float64   `json:"long" gorm:"not null"`
	RoomingHouseID uuid.UUID `json:"rooming_house_id" gorm:"not null"`
	Rooms          []Room    `json:"rooms" gorm:"foreignKey:SizeID"`
}

type AddSizeBody struct {
	Name           string    `json:"name" validate:"required"`
	Width          float64   `json:"width" validate:"required"`
	Long           float64   `json:"long" validate:"required"`
	RoomingHouseID uuid.UUID `json:"rooming_house_id"`
}

type UpdateSizeBody struct {
	Name  string  `json:"name" validate:"required"`
	Width float64 `json:"width" validate:"required"`
	Long  float64 `json:"long" validate:"required"`
}

type AllSizeResponse struct {
	ID           uuid.UUID                  `json:"id"`
	Name         string                     `json:"name"`
	Width        float64                    `json:"width"`
	Long         float64                    `json:"long"`
	RoomingHouse TenantRoomingHouseResponse `json:"rooming_house"`
}

func (s *Size) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	s.CreatedAt = time.Now()

	return
}
