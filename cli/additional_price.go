package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/middlewares"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func AdditionalPriceRoutes(e *echo.Echo) {
	additionalPriceRepo := repositories.NewAdditionalPriceRepository(config.DB)
	additionalPeriodRepo := repositories.NewAdditionalPeriodRepository(config.DB)
	periodRepo := repositories.NewPeriodRepository(config.DB)
	roomingHouseRepo := repositories.NewRoomingHouseRepository(config.DB)

	additionalPriceController := controllers.NewAdditionalPriceController(additionalPriceRepo, additionalPeriodRepo, periodRepo, roomingHouseRepo)

	additionalPrice := e.Group("/additionals")
	additionalPrice.GET("/:id", additionalPriceController.FindAdditionalPriceByID, middlewares.JWTAuth)
	additionalPrice.GET("", additionalPriceController.FindAllAdditionalPrices, middlewares.JWTAuth)
	additionalPrice.POST("", additionalPriceController.CreateAdditionalPrice, middlewares.JWTAuth, middlewares.Authz)
	additionalPrice.PUT("/:id", additionalPriceController.UpdateAdditionalPriceByID, middlewares.JWTAuth, middlewares.Authz)
	additionalPrice.DELETE("/:id", additionalPriceController.DeleteAdditionalPriceByID, middlewares.JWTAuth, middlewares.Authz)
}
