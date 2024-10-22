package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	BaseModel
	Day                   int        `json:"day" gorm:"not null"`
	Month                 int        `json:"month" gorm:"not null"`
	Year                  int        `json:"year" gorm:"not null"`
	Amount                float64    `json:"amount" gorm:"not null"`
	Description           string     `json:"description" gorm:"not null"`
	IsRoom                bool       `json:"is_room" gorm:"not null"`
	TransactionCategoryID uuid.UUID  `json:"transaction_category_id" gorm:"not null;size:191"`
	RoomID                *uuid.UUID `json:"room_id" gorm:"size:191"`
	TenantID              *uuid.UUID `json:"tenant_id" gorm:"size:191"`
	RoomingHouseID        uuid.UUID  `json:"rooming_house_id" gorm:"not null;size:191"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()

	return
}
