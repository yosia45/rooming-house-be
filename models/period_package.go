package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeriodPackage struct {
	BaseModel
	PeriodID         uuid.UUID `json:"period_id" gorm:"not null;size:191"`
	PricingPackageID uuid.UUID `json:"pricing_package_id" gorm:"not null;size:191"`
	Price            float64   `json:"price" gorm:"not null"`
	Period           Period    `json:"period" gorm:"foreignKey:PeriodID"`
}

func (pp *PeriodPackage) BeforeCreate(tx *gorm.DB) (err error) {
	pp.ID = uuid.New()
	location, _ := time.LoadLocation("Asia/Jakarta")
	pp.CreatedAt = time.Now().In(location)

	return
}
