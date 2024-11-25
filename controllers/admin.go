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

func (ac *AdminController) DeleteAdminByID(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)
	adminID := c.Param("id")

	if userPayload.Role != "owner" {
		return utils.HandlerError(c, utils.NewForbiddenError("you are not allowed to access this resource"))
	}

	adminUUID, err := uuid.Parse(adminID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid admin id"))
	}

	if err := ac.adminRepo.DeleteAdminByID(adminUUID); err != nil {
		return utils.HandlerError(c, utils.NewInternalError(err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "admin deleted successfully"})
}
