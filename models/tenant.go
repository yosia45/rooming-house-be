package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	BaseModel
	Name                   string                  `json:"name" gorm:"not null"`
	Gender                 string                  `json:"gender" gorm:"not null"`
	PhoneNumber            string                  `json:"phoneNumber" gorm:"not null"`
	EmergencyContact       string                  `json:"emergencyContact" gorm:"not null"`
	StartDate              time.Time               `json:"start_date" gorm:"not null"`
	Duration               int                     `json:"duration" gorm:"not null"`
	IsTenant               bool                    `json:"is_tenant" gorm:"not null"`
	IsDepositPaid          bool                    `json:"is_deposit_paid" gorm:"not null"`
	IsDepositBack          bool                    `json:"is_deposit_back" gorm:"not null"`
	RoomingHouseID         uuid.UUID               `json:"rooming_house_id" gorm:"not null"`
	PeriodID               uuid.UUID               `json:"period_id" gorm:"not null;size:191"`
	TenantID               uuid.UUID               `json:"tenant_id" gorm:"size:191"`
	RoomID                 uuid.UUID               `json:"room_id" gorm:"not null;size:191"`
	TenantAdditionalPrices []TenantAdditionalPrice `json:"tenant_additional_price" gorm:"foreignKey:TenantID"`
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()

	return
}
