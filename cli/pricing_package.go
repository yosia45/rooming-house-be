package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/middlewares"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func PricingPackageRoutes(e *echo.Echo) {
	pricingPackageRepo := repositories.NewPricingPackageRepository(config.DB)
	periodRepo := repositories.NewPeriodRepository(config.DB)
	periodPackageRepo := repositories.NewPeriodPackageRepository(config.DB)
	roomingHouseRepo := repositories.NewRoomingHouseRepository(config.DB)

	pricingPackageController := controllers.NewPricingPackageController(pricingPackageRepo, periodRepo, periodPackageRepo, roomingHouseRepo)

	pricingPackage := e.Group("/packages")
	pricingPackage.POST("", pricingPackageController.CreatePricingPackage, middlewares.JWTAuth, middlewares.Authz)
	pricingPackage.GET("", pricingPackageController.GetAllPricingPackages, middlewares.JWTAuth)
	pricingPackage.PUT("/:id", pricingPackageController.UpdatePricingPackage, middlewares.JWTAuth, middlewares.Authz)
	pricingPackage.DELETE("/:id", pricingPackageController.DeletePricingPackage, middlewares.JWTAuth, middlewares.Authz)
}
