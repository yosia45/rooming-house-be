package main

import (
	"log"
	"os"
	"rooming-house-cms-be/cli"
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/middlewares"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()

	// seeders.SeedPeriod(config.DB)
	// seeders.SeedFacility(config.DB)
	// seeders.SeedTransactionCategory(config.DB)

	port := os.Getenv("PORT")

	e := echo.New()

	e.Use(middlewares.CORSConfig())

	cli.Auth(e)
	cli.RoomingHouseRoutes(e)
	cli.SizeRoutes(e)
	cli.AdditionalPriceRoutes(e)
	cli.TransactionCategoryRoutes(e)
	cli.PricingPackageRoutes(e)
	cli.RoomRoutes(e)
	cli.AdditionalPriceRoutes(e)
	cli.TransactionRoutes(e)
	cli.TenantRoutes(e)
	cli.AdminRoutes(e)
	cli.PeriodRoute(e)

	e.Logger.Fatal(e.Start(":" + port))
}
