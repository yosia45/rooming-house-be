package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func PricingPackageRoutes(e *echo.Echo) {
	pricingPackageRepo := repositories.NewPricingPackageRepository(config.DB)
	periodRepo := repositories.NewPeriodRepository(config.DB)
	periodPackageRepo := repositories.NewPeriodPackageRepository(config.DB)

	pricingPackageController := controllers.NewPricingPackageController(pricingPackageRepo, periodRepo, periodPackageRepo)

	pricingPackage := e.Group("/pricing-package")
	pricingPackage.POST("", pricingPackageController.CreatePricingPackage)
	pricingPackage.GET("", pricingPackageController.GetAllPricingPackages)
	pricingPackage.PUT("/:id", pricingPackageController.UpdatePricingPackage)
	pricingPackage.DELETE("/:id", pricingPackageController.DeletePricingPackage)
}
