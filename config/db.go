package config

import (
	"fmt"
	"log"
	"os"
	"rooming-house-cms-be/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", username, password, host, port, name)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect DB: ", err)
	}

	DB.AutoMigrate(
		&models.Owner{},
		&models.RoomingHouse{},
		&models.Admin{},
		&models.PricingPackage{},
		&models.Period{},
		&models.PeriodPackage{},
		&models.Facility{},
		&models.RoomingHouseFacility{},
		&models.Size{},
		&models.Room{},
		&models.RoomFacility{},
		&models.TransactionCategory{},
		&models.Transaction{},
		&models.AdditionalPrice{},
		&models.AdditionalPeriod{},
		&models.Tenant{},
		&models.TenantAdditionalPrice{},
	)

	log.Println("Success connecting to DB")
}
