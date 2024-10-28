package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func TenantRoutes(e *echo.Echo) {
	tenantRepo := repositories.NewTenantRepository(config.DB)
	tenantAdditionalRepo := repositories.NewTenantAdditionalRepository(config.DB)

	tenantController := controllers.NewTenantController(tenantRepo, tenantAdditionalRepo)

	tenant := e.Group("/tenant")
	tenant.POST("", tenantController.CreateTenant)
	tenant.GET("", tenantController.FindAllTenants)
	tenant.GET("/:id", tenantController.FindTenantByID)
	tenant.DELETE("/:id", tenantController.DeleteTenantByID)
}
