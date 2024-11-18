package controllers

import (
	"net/http"
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AdminController struct {
	adminRepo        repositories.AdminRepository
	roomingHouseRepo repositories.RoomingHouseRepository
}

func NewAdminController(adminRepo repositories.AdminRepository, roomingHouseRepo repositories.RoomingHouseRepository) *AdminController {
	return &AdminController{adminRepo: adminRepo, roomingHouseRepo: roomingHouseRepo}
}

func (ac *AdminController) GetAllAdmin(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)

	var roomingHouseIDs []uuid.UUID

	if userPayload.RoomingHouseID == uuid.Nil {
		roomingHouses, err := ac.roomingHouseRepo.FindAllRoomingHouse(userPayload.RoomingHouseID, userPayload.UserID, userPayload.Role)
		if err != nil {
			return utils.HandlerError(c, utils.NewInternalError(err.Error()))
		}

		for _, roomingHouse := range roomingHouses {
			roomingHouseIDs = append(roomingHouseIDs, roomingHouse.ID)
		}
	}

	admins, err := ac.adminRepo.FindAllAdmin(roomingHouseIDs)
	if err != nil {
		return utils.HandlerError(c, utils.NewInternalError(err.Error()))
	}

	return c.JSON(http.StatusOK, admins)
}
