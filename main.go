package main

import (
	"log"
	"os"
	"rooming-house-cms-be/cli"
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/seeders"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()

	seeders.SeedPeriod(config.DB)
	seeders.SeedFacility(config.DB)
	seeders.SeedTransactionCategory(config.DB)

	port := os.Getenv("PORT")

	e := echo.New()

	cli.Auth(e)
	cli.RoomingHouseRoutes(e)
	cli.SizeRoutes(e)
	cli.AdditionalPriceRoutes(e)
	cli.TransactionCategoryRoutes(e)

	e.Logger.Fatal(e.Start(":" + port))
}
