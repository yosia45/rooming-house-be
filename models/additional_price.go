package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdditionalPrice struct {
	BaseModel
	Name                   string                  `json:"name" gorm:"not null"`
	Price                  float64                 `json:"price" gorm:"not null"`
	RoomingHouseID         uuid.UUID               `json:"rooming_house_id" gorm:"not null;size:191"`
	TenantAdditionalPrices []TenantAdditionalPrice `json:"tenant_additional_prices" gorm:"foreignKey:AdditionalPriceID"`
}

type AddAdditionalPriceBody struct {
	Name           string    `json:"name"`
	Price          float64   `json:"price"`
	RoomingHouseID uuid.UUID `json:"rooming_house_id"`
}

type UpdateAdditionalPriceBody struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (ap *AdditionalPrice) BeforeCreate(tx *gorm.DB) (err error) {
	ap.ID = uuid.New()
	ap.CreatedAt = time.Now()

	return
}
