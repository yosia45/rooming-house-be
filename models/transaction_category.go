package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionCategory struct {
	BaseModel
	Name         string        `json:"name" gorm:"not null"`
	IsExpense    bool          `json:"is_expense" gorm:"not null"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:TransactionCategoryID"`
}

type TransactionCategoryBody struct {
	Name      string `json:"name" gorm:"column:transaction_category_name"`
	IsExpense bool   `json:"is_expense" gorm:"column:transaction_category_is_expense"`
}

func (tc *TransactionCategory) BeforeCreate(tx *gorm.DB) (err error) {
	tc.ID = uuid.New()
	tc.CreatedAt = time.Now()

	return
}
