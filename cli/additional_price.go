package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func AdditionalPriceRoutes(e *echo.Echo) {
	additionalPriceRepo := repositories.NewAdditionalPriceRepository(config.DB)

	additionalPriceController := controllers.NewAdditionalPriceController(additionalPriceRepo)

	additionalPrice := e.Group("/additional-price")
	additionalPrice.GET("/:id", additionalPriceController.FindAdditionalPriceByID)
	additionalPrice.GET("", additionalPriceController.FindAllAdditionalPrices)
	additionalPrice.POST("", additionalPriceController.CreateAdditionalPrice)
	additionalPrice.PUT("/:id", additionalPriceController.UpdateAdditionalPriceByID)
	additionalPrice.DELETE("/:id", additionalPriceController.DeleteAdditionalPriceByID)
}
