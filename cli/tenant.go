package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/middlewares"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func TenantRoutes(e *echo.Echo) {
	tenantRepo := repositories.NewTenantRepository(config.DB)
	tenantAdditionalRepo := repositories.NewTenantAdditionalRepository(config.DB)
	roomingHouseRepo := repositories.NewRoomingHouseRepository(config.DB)
	roomRepo := repositories.NewRoomRepository(config.DB)

	tenantController := controllers.NewTenantController(tenantRepo, tenantAdditionalRepo, roomingHouseRepo, roomRepo)

	tenant := e.Group("/tenants", middlewares.JWTAuth)
	tenant.POST("", tenantController.CreateTenant)
	tenant.GET("", tenantController.FindAllTenants)
	tenant.GET("/:id", tenantController.FindTenantByID)
	tenant.DELETE("/:id", tenantController.DeleteTenantByID)
}
