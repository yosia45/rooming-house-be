package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func FacilityRoutes(e *echo.Echo) {
	facilityRepo := repositories.NewFacilityRepository(config.DB)

	facilityController := controllers.NewFacilityController(facilityRepo)

	facility := e.Group("/facilities")
	facility.GET("", facilityController.GetAllFacilities)
}
