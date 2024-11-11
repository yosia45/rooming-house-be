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
	Description           string     `json:"description"`
	IsRoom                bool       `json:"is_room" gorm:"not null"`
	TransactionCategoryID uuid.UUID  `json:"transaction_category_id" gorm:"not null;size:191"`
	RoomID                *uuid.UUID `json:"room_id" gorm:"size:191"`
	TenantID              *uuid.UUID `json:"tenant_id" gorm:"size:191"`
	RoomingHouseID        uuid.UUID  `json:"rooming_house_id" gorm:"not null;size:191"`
}

type AddTransactionBody struct {
	Day                   int        `json:"day"`
	Month                 int        `json:"month"`
	Year                  int        `json:"year"`
	IsRoom                bool       `json:"is_room"`
	Description           string     `json:"description"`
	Amount                float64    `json:"amount"`
	TransactionCategoryID uuid.UUID  `json:"transaction_category_id"`
	RoomID                *uuid.UUID `json:"room_id"`
	TenantID              *uuid.UUID `json:"tenant_id"`
	RoomingHouseID        uuid.UUID  `json:"rooming_house_id"`
}

type TransactionResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Day         int       `json:"day"`
	Month       int       `json:"month"`
	Year        int       `json:"year"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	IsExpense   bool      `json:"is_expense"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()

	return
}
