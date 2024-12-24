package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/middlewares"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func PeriodRoute(e *echo.Echo) {
	periodRepo := repositories.NewPeriodRepository(config.DB)

	periodController := controllers.NewPeriodController(periodRepo)

	period := e.Group("/periods")
	period.GET("", periodController.GetAllPeriods, middlewares.JWTAuth)
}
