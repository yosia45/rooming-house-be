package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	BaseModel
	Name                   string            `json:"name" gorm:"not null"`
	Gender                 string            `json:"gender" gorm:"not null"`
	PhoneNumber            string            `json:"phoneNumber" gorm:"not null"`
	EmergencyContact       string            `json:"emergencyContact"`
	StartDate              *time.Time        `json:"start_date"`
	EndDate                *time.Time        `json:"end_date"`
	RegularPaymentDuration int               `json:"regular_payment_duration"`
	IsTenant               bool              `json:"is_tenant" gorm:"not null"`
	IsDepositPaid          bool              `json:"is_deposit_paid"`
	IsDepositBack          bool              `json:"is_deposit_back"`
	RoomingHouseID         uuid.UUID         `json:"rooming_house_id" gorm:"not null"`
	PeriodID               *uuid.UUID        `json:"period_id" gorm:"size:191"`
	TenantID               uuid.UUID         `json:"tenant_id" gorm:"size:191"`
	RoomID                 *uuid.UUID        `json:"room_id" gorm:"size:191"`
	Transactions           []Transaction     `json:"transactions" gorm:"foreignKey:TenantID;references:ID"`
	AdditionalPrices       []AdditionalPrice `json:"additional_prices" gorm:"many2many:tenant_additional_prices;joinForeignKey:TenantID;joinReferences:AdditionalPriceID"`
}

type AddTenantBody struct {
	Name                   string      `json:"name"`
	Gender                 string      `json:"gender"`
	PhoneNumber            string      `json:"phoneNumber"`
	EmergencyContact       string      `json:"emergencyContact"`
	IsTenant               bool        `json:"is_tenant"`
	RegularPaymentDuration int         `json:"regular_payment_duration"`
	RoomingHouseID         uuid.UUID   `json:"rooming_house_id"`
	RoomID                 *uuid.UUID  `json:"room_id"`
	PeriodID               *uuid.UUID  `json:"period_id"`
	TenantID               uuid.UUID   `json:"tenant_id"`
	TenantAdditionalIDs    []uuid.UUID `json:"tenant_additional_ids"`
}

type GetAllTenantResponse struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	Gender         string     `json:"gender"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	RoomID         uuid.UUID  `json:"room_id"`
	RoomingHouseID uuid.UUID  `json:"rooming_house_id"`
}

type AllTenantRepoResponse struct {
	ID           uuid.UUID                  `json:"id"`
	Name         string                     `json:"name"`
	Gender       string                     `json:"gender"`
	IsTenant     bool                       `json:"is_tenant"`
	StartDate    *time.Time                 `json:"start_date"`
	EndDate      *time.Time                 `json:"end_date"`
	Room         *TenantRoomResponse        `json:"room" gorm:"embedded"`
	RoomingHouse TenantRoomingHouseResponse `json:"rooming_house" gorm:"embedded"`
}

type TenantDetailResponse struct {
	ID                     uuid.UUID                  `json:"id"`
	Name                   string                     `json:"name"`
	Gender                 string                     `json:"gender"`
	PhoneNumber            string                     `json:"phoneNumber"`
	EmergencyContact       string                     `json:"emergencyContact"`
	BookedRoomID           uuid.UUID                  `json:"booked_room_id"`
	StartDate              *time.Time                 `json:"start_date"`
	EndDate                *time.Time                 `json:"end_date"`
	RegularPaymentDuration int                        `json:"regular_payment_duration"`
	IsTenant               bool                       `json:"is_tenant"`
	IsDepositPaid          bool                       `json:"is_deposit_paid"`
	IsDepositBack          bool                       `json:"is_deposit_back"`
	TenantID               uuid.UUID                  `json:"tenant_id"`
	RoomingHouse           TenantRoomingHouseResponse `json:"rooming_house" gorm:"embedded"`
	Period                 PeriodResponse             `json:"period" gorm:"embedded"`
	TenantAssists          []TenantAssistResponse     `json:"tenant_assists" gorm:"-"`
	Room                   TenantRoomResponse         `json:"room" gorm:"embedded"`
	Transactions           []TransactionResponse      `json:"transactions" gorm:"-"`
	AdditionalPrices       []AdditionalPriceDetail    `json:"additional_prices" gorm:"-"`
}

type TenantRoomDetailResponse struct {
	ID                     uuid.UUID              `json:"id"`
	Name                   string                 `json:"name"`
	Gender                 string                 `json:"gender"`
	PhoneNumber            string                 `json:"phoneNumber"`
	EmergencyContact       string                 `json:"emergencyContact"`
	IsTenant               bool                   `json:"is_tenant"`
	StartDate              *time.Time             `json:"start_date"`
	EndDate                *time.Time             `json:"end_date"`
	RegularPaymentDuration int                    `json:"regular_payment_duration"`
	IsDepositPaid          bool                   `json:"is_deposit_paid"`
	IsDepositBack          bool                   `json:"is_deposit_back"`
	TenantAssists          []TenantAssistResponse `json:"tenant_assists" gorm:"-"`
}

type TenantAssistResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Gender         string    `json:"gender"`
	PhoneNumber    string    `json:"phoneNumber"`
	IsTenant       bool      `json:"is_tenant"`
	TenantID       uuid.UUID `json:"tenant_id"`
	RoomingHouseID uuid.UUID `json:"rooming_house_id"`
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()

	return
}
