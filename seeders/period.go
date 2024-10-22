package seeders

import (
	"rooming-house-cms-be/models"

	"gorm.io/gorm"
)

func SeedPeriod(db *gorm.DB) {
	periods := []models.Period{
		{
			Name: "Daily",
			Unit: "day",
		},
		{
			Name: "Weekly",
			Unit: "week",
		},
		{
			Name: "Monthly",
			Unit: "month",
		},
		{
			Name: "Yearly",
			Unit: "year",
		},
	}

	for _, period := range periods {
		db.Create(&period)
	}
}
