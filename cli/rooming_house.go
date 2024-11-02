package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/middlewares"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func RoomingHouseRoutes(e *echo.Echo) {
	roomingHouseRepo := repositories.NewRoomingHouseRepository(config.DB)
	roomingHouseFacilityRepo := repositories.NewRoomingHouseFacilityRepository(config.DB)
	facilityRepo := repositories.NewFacilityRepository(config.DB)

	roomingHouseController := controllers.NewRoomingHouseController(roomingHouseRepo, roomingHouseFacilityRepo, facilityRepo)

	roomingHouse := e.Group("/roominghouses")
	roomingHouse.GET("/:id", roomingHouseController.GetRoomingHouseByID, middlewares.JWTAuth)
	roomingHouse.GET("", roomingHouseController.GetAllRoomingHouse, middlewares.JWTAuth)
	roomingHouse.POST("", roomingHouseController.CreateRoomingHouse, middlewares.JWTAuth, middlewares.Authz)
	roomingHouse.PUT("/:id", roomingHouseController.UpdateRoomingHouseByID, middlewares.JWTAuth, middlewares.Authz)
	roomingHouse.DELETE("/:id", roomingHouseController.DeleteRoomingHouseByID, middlewares.JWTAuth, middlewares.Authz)
}
