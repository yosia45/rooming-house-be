package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PricingPackage struct {
	BaseModel
	Name           string          `json:"name" gorm:"not null"`
	RoomingHouseID uuid.UUID       `json:"rooming_house_id" gorm:"not null"`
	PeriodPackages []PeriodPackage `json:"period_packages" gorm:"foreignKey:PricingPackageID"`
	Rooms          []Room          `json:"rooms" gorm:"foreignKey:PackageID"`
}

type AddPricingPackageBody struct {
	Name           string    `json:"name"`
	RoomingHouseID uuid.UUID `json:"rooming_house_id"`
	DailyPrice     float64   `json:"daily_price"`
	WeeklyPrice    float64   `json:"weekly_price"`
	MonthlyPrice   float64   `json:"monthly_price"`
	AnnualPrice    float64   `json:"annual_price"`
}

type UpdatePricingPackageBody struct {
	Name         string  `json:"name"`
	DailyPrice   float64 `json:"daily_price"`
	WeeklyPrice  float64 `json:"weekly_price"`
	MonthlyPrice float64 `json:"monthly_price"`
	AnnualPrice  float64 `json:"annual_price"`
}

func (pp *PricingPackage) BeforeCreate(tx *gorm.DB) (err error) {
	pp.ID = uuid.New()
	pp.CreatedAt = time.Now()

	return
}
