package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdditionalPeriod struct {
	BaseModel
	AdditionalPriceID uuid.UUID `json:"additional_price_id" gorm:"not null;size:191"`
	PeriodID          uuid.UUID `json:"period_id" gorm:"not null;size:191"`
	Price             float64   `json:"price" gorm:"not null"`
	Period            Period    `json:"period" gorm:"foreignKey:PeriodID"`
}

type AddAdditionalPeriodBody struct {
	AdditionalPriceID uuid.UUID `json:"additional_price_id"`
	PeriodID          uuid.UUID `json:"period_id"`
	Price             float64   `json:"price"`
}

func (ap *AdditionalPeriod) BeforeCreate(tx *gorm.DB) (err error) {
	ap.ID = uuid.New()
	ap.CreatedAt = time.Now()

	return
}
