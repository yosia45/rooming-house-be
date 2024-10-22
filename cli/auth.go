package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/middlewares"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func Auth(e *echo.Echo) {
	adminRepo := repositories.NewAdminRepository(config.DB)
	ownerRepo := repositories.NewOwnerRepository(config.DB)
	userController := controllers.NewUserController(ownerRepo, adminRepo)

	e.POST("/login", userController.Login)
	e.POST("/registerowner", userController.RegisterOwner)
	e.POST("/registeradmin", userController.RegisterAdmin, middlewares.JWTAuth, middlewares.Authz)
}
