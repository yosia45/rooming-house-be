package seeders

import (
	"rooming-house-cms-be/models"

	"gorm.io/gorm"
)

func SeedFacility(db *gorm.DB) {
	facilities := []models.Facility{
		{
			Name:     "AC",
			IsPublic: false,
			IsRoom:   true,
		},
		{
			Name:     "TV",
			IsPublic: false,
			IsRoom:   true,
		},
		{
			Name:     "WiFi",
			IsPublic: true,
			IsRoom:   false,
		},
		{
			Name:     "Bathroom",
			IsPublic: false,
			IsRoom:   true,
		},
		{
			Name:     "Parking",
			IsPublic: true,
			IsRoom:   false,
		},
		{
			Name:     "Kitchen",
			IsPublic: true,
			IsRoom:   false,
		},
		{
			Name:     "Wardrobe",
			IsPublic: false,
			IsRoom:   true,
		},
		{
			Name:     "Single-Size Bed",
			IsPublic: false,
			IsRoom:   true,
		},
	}

	for _, facility := range facilities {
		db.Create(&facility)
	}
}
