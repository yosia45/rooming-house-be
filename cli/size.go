package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/middlewares"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func SizeRoutes(e *echo.Echo) {
	sizeRepo := repositories.NewSizeRepository(config.DB)

	sizeController := controllers.NewSizeController(sizeRepo)

	size := e.Group("/size")
	size.GET("", sizeController.FindAllSizes, middlewares.JWTAuth)
	size.GET("/:id", sizeController.FindSizeByID, middlewares.JWTAuth)
	size.POST("", sizeController.CreateSize, middlewares.JWTAuth, middlewares.Authz)
	size.PUT("/:id", sizeController.UpdateSizeByID, middlewares.JWTAuth, middlewares.Authz)
	size.DELETE("/:id", sizeController.DeleteSizeByID, middlewares.JWTAuth, middlewares.Authz)
}
