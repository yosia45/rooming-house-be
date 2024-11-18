package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/middlewares"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func AdminRoutes(e *echo.Echo) {
	adminRepo := repositories.NewAdminRepository(config.DB)
	roomingHouseRepo := repositories.NewRoomingHouseRepository(config.DB)

	adminController := controllers.NewAdminController(adminRepo, roomingHouseRepo)

	admin := e.Group("/admins")
	admin.GET("", adminController.GetAllAdmin, middlewares.JWTAuth, middlewares.Authz)
}
