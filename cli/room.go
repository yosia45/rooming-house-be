package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/middlewares"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func RoomRoutes(e *echo.Echo) {
	roomRepo := repositories.NewRoomRepository(config.DB)
	roomFacilityRepo := repositories.NewRoomFacilityRepository(config.DB)
	roomingHouseRepo := repositories.NewRoomingHouseRepository(config.DB)
	sizeRepo := repositories.NewSizeRepository(config.DB)
	packageRepo := repositories.NewPricingPackageRepository(config.DB)
	facilityRepo := repositories.NewFacilityRepository(config.DB)

	roomController := controllers.NewRoomController(roomRepo, roomFacilityRepo, roomingHouseRepo, sizeRepo, packageRepo, facilityRepo)

	room := e.Group("/rooms")
	room.POST("", roomController.CreateRoom, middlewares.JWTAuth)
	room.GET("", roomController.GetAllRooms, middlewares.JWTAuth)
	room.GET("/:id", roomController.GetRoomByID, middlewares.JWTAuth)
	room.PUT("/:id", roomController.UpdateRoomByID, middlewares.JWTAuth)
	room.DELETE("/:id", roomController.DeleteRoomByID, middlewares.JWTAuth, middlewares.Authz)
}
