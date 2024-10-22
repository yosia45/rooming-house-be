package controllers

import (
	"fmt"
	"net/http"
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RoomingHouseController struct {
	roomingHouseRepo         repositories.RoomingHouseRepository
	roomingHouseFacilityRepo repositories.RoomingHouseFacilityRepository
}

func NewRoomingHouseController(roomingHouseRepo repositories.RoomingHouseRepository, roomingHouseFacilityRepo repositories.RoomingHouseFacilityRepository) *RoomingHouseController {
	return &RoomingHouseController{roomingHouseRepo: roomingHouseRepo, roomingHouseFacilityRepo: roomingHouseFacilityRepo}
}

func (rhc *RoomingHouseController) CreateRoomingHouse(c echo.Context) error {
	fmt.Println("masuk")
	userPayload := c.Get("userPayload").(*models.JWTPayload)

	var roomingHouseBody models.RoomingHouseBody
	if err := c.Bind(&roomingHouseBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid request body"))
	}

	if roomingHouseBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	if roomingHouseBody.Address == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("address is required"))
	}

	if roomingHouseBody.Description == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("description is required"))
	}

	if roomingHouseBody.FloorTotal == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("floor total is required"))
	}

	if len(roomingHouseBody.RoomingHouseFacilityIDs) == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("rooming house facility is required"))
	}

	newRoomingHouse := models.RoomingHouse{
		Name:        roomingHouseBody.Name,
		Address:     roomingHouseBody.Address,
		Description: roomingHouseBody.Description,
		FloorTotal:  roomingHouseBody.FloorTotal,
		OwnerID:     userPayload.UserID,
	}

	if err := rhc.roomingHouseRepo.CreateRoomingHouse(&newRoomingHouse); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to create rooming house"))
	}

	var RoomingHouseFacilities []models.RoomingHouseFacility
	for _, roomingHouseFacilityID := range roomingHouseBody.RoomingHouseFacilityIDs {
		roomingHouseFacility := models.RoomingHouseFacility{
			RoomingHouseID: newRoomingHouse.ID,
			FacilityID:     roomingHouseFacilityID,
		}
		RoomingHouseFacilities = append(RoomingHouseFacilities, roomingHouseFacility)
	}

	if err := rhc.roomingHouseFacilityRepo.CreateRoomingHouseFacility(&RoomingHouseFacilities); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to create rooming house facility"))
	}

	return c.JSON(http.StatusCreated, newRoomingHouse)
}

func (rhc *RoomingHouseController) GetAllRoomingHouse(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)
	roomingHouses, err := rhc.roomingHouseRepo.FindAllRoomingHouse(userPayload.RoomingHouseID, userPayload.UserID, userPayload.Role)
	if err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to get all rooming house"))
	}

	return c.JSON(http.StatusOK, roomingHouses)
}

func (rhc *RoomingHouseController) GetRoomingHouseByID(c echo.Context) error {
	roomingHouseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid rooming house ID"))
	}

	roomingHouse, err := rhc.roomingHouseRepo.FindRoomingHouseByID(roomingHouseID)
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("rooming house not found"))
	}

	return c.JSON(http.StatusOK, roomingHouse)
}

func (rhc *RoomingHouseController) UpdateRoomingHouseByID(c echo.Context) error {
	roomingHouseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid rooming house ID"))
	}

	var roomingHouseBody models.RoomingHouseBody
	if err := c.Bind(&roomingHouseBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid request body"))
	}

	if roomingHouseBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	if roomingHouseBody.Address == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("address is required"))
	}

	if roomingHouseBody.Description == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("description is required"))
	}

	if roomingHouseBody.FloorTotal == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("floor total is required"))
	}

	if len(roomingHouseBody.RoomingHouseFacilityIDs) == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("rooming house facility is required"))
	}

	edittedRoomingHouse := models.RoomingHouse{
		Name:        roomingHouseBody.Name,
		Address:     roomingHouseBody.Address,
		Description: roomingHouseBody.Description,
		FloorTotal:  roomingHouseBody.FloorTotal,
	}

	if err := rhc.roomingHouseRepo.UpdateRoomingHouse(&edittedRoomingHouse, roomingHouseID); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to update rooming house"))
	}

	var RoomingHouseFacilities []models.RoomingHouseFacility
	for _, roomingHouseFacilityID := range roomingHouseBody.RoomingHouseFacilityIDs {
		roomingHouseFacility := models.RoomingHouseFacility{
			RoomingHouseID: roomingHouseID,
			FacilityID:     roomingHouseFacilityID,
		}
		RoomingHouseFacilities = append(RoomingHouseFacilities, roomingHouseFacility)
	}

	if err := rhc.roomingHouseFacilityRepo.UpdateRoomingHouseFacilityByRoomingHouseID(&RoomingHouseFacilities, roomingHouseID); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to update rooming house facility"))
	}

	return c.JSON(http.StatusOK, edittedRoomingHouse)
}

func (rhc *RoomingHouseController) DeleteRoomingHouseByID(c echo.Context) error {
	roomingHouseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid rooming house ID"))
	}

	if err := rhc.roomingHouseRepo.DeleteRoomingHouse(roomingHouseID); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to delete rooming house"))
	}

	return c.NoContent(http.StatusNoContent)
}
