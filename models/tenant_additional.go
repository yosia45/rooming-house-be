package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantAdditionalPrice struct {
	BaseModel
	TenantID          uuid.UUID `json:"tenant_id" gorm:"not null;size:191"`
	AdditionalPriceID uuid.UUID `json:"additional_price_id" gorm:"not null;size:191"`
}

type AddTenantAdditionalPriceBody struct {
	TenantID          uuid.UUID `json:"tenant_id"`
	AdditionalPriceID uuid.UUID `json:"additional_price_id"`
}

func (tap *TenantAdditionalPrice) BeforeCreate(tx *gorm.DB) (err error) {
	tap.ID = uuid.New()
	tap.CreatedAt = time.Now()

	return
}
