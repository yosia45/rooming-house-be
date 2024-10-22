package seeders

import (
	"rooming-house-cms-be/models"

	"gorm.io/gorm"
)

func SeedTransactionCategory(db *gorm.DB) {
	transactionCategories := []models.TransactionCategory{
		{
			Name:      "Deposit",
			IsExpense: false,
		},
		{
			Name:      "Rent",
			IsExpense: false,
		},
		{
			Name:      "Maintenance",
			IsExpense: true,
		},
		{
			Name:      "Salary",
			IsExpense: true,
		},
		{
			Name:      "Repairment",
			IsExpense: true,
		},
		{
			Name:      "Utilities",
			IsExpense: true,
		},
		{
			Name:      "Deposit Payback",
			IsExpense: true,
		},
	}

	for _, transactionCategory := range transactionCategories {
		db.Create(&transactionCategory)
	}
}
