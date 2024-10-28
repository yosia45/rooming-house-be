package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdditionalPrice struct {
	BaseModel
	Name                   string                  `json:"name" gorm:"not null"`
	RoomingHouseID         uuid.UUID               `json:"rooming_house_id" gorm:"not null;size:191"`
	TenantAdditionalPrices []TenantAdditionalPrice `json:"tenant_additional_prices" gorm:"foreignKey:AdditionalPriceID"`
	AdditionalPeriods      []AdditionalPeriod      `json:"additional_periods" gorm:"foreignKey:AdditionalPriceID"`
}

type AddAdditionalPriceBody struct {
	Name           string    `json:"name"`
	DailyPrice     float64   `json:"daily_price"`
	WeeklyPrice    float64   `json:"weekly_price"`
	MonthlyPrice   float64   `json:"monthly_price"`
	AnnualPrice    float64   `json:"annual_price"`
	RoomingHouseID uuid.UUID `json:"rooming_house_id"`
}

type UpdateAdditionalPriceBody struct {
	Name         string  `json:"name"`
	DailyPrice   float64 `json:"daily_price"`
	WeeklyPrice  float64 `json:"weekly_price"`
	MonthlyPrice float64 `json:"monthly_price"`
	AnnualPrice  float64 `json:"annual_price"`
}

func (ap *AdditionalPrice) BeforeCreate(tx *gorm.DB) (err error) {
	ap.ID = uuid.New()
	ap.CreatedAt = time.Now()

	return
}
