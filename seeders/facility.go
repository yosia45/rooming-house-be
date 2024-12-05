package seeders

import (
	"rooming-house-cms-be/models"

	"gorm.io/gorm"
)

func SeedFacility(db *gorm.DB) {
	facilities := []models.Facility{
		{
			Name:     "Air  Conditioner",
			IsPublic: false,
			IsRoom:   true,
		},
		{
			Name:        "TV",
			Description: "43 Inch",
			IsPublic:    false,
			IsRoom:      true,
		},
		{
			Name:     "Wi-Fi",
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
		{
			Name:     "Double-Size Bed",
			IsPublic: false,
			IsRoom:   true,
		},
		{
			Name:        "Water Heater",
			Description: "Bathroom Facility",
			IsPublic:    false,
			IsRoom:      true,
		},
		{
			Name:     "Fridge",
			IsPublic: true,
			IsRoom:   false,
		},
		{
			Name:     "Working Desk and Chair Set",
			IsPublic: false,
			IsRoom:   true,
		},
		{
			Name:     "Dispenser",
			IsPublic: true,
			IsRoom:   false,
		},
		{
			Name:     "Balcony",
			IsPublic: false,
			IsRoom:   true,
		},
	}

	for _, facility := range facilities {
		db.Create(&facility)
	}
}
