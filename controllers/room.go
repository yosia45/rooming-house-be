package controllers

import (
	"net/http"
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RoomController struct {
	roomRepo         repositories.RoomRepository
	roomFacilityRepo repositories.RoomFacilityRepository
}

func NewRoomController(roomRepo repositories.RoomRepository, roomFacilityRepo repositories.RoomFacilityRepository) *RoomController {
	return &RoomController{roomRepo: roomRepo, roomFacilityRepo: roomFacilityRepo}
}

func (rc *RoomController) CreateRoom(c echo.Context) error {
	var roomBody models.AddRoomBody
	userPayload := c.Get("userPayload").(*models.JWTPayload)

	if err := c.Bind(&roomBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if roomBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	if roomBody.Floor == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("floor is required"))
	}

	if roomBody.MaxCapacity == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("max capacity is required"))
	}

	if roomBody.SizeID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("size id is required"))
	}

	if roomBody.PackageID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("pricing id is required"))
	}

	if len(roomBody.RoomFacilities) == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("facilities is required"))
	}

	if roomBody.RoomingHouseID == uuid.Nil && userPayload.Role == "owner" {
		return utils.HandlerError(c, utils.NewBadRequestError("rooming house id is required"))
	}

	var roomingHouseID uuid.UUID

	if userPayload.Role == "owner" {
		roomingHouseID = roomBody.RoomingHouseID
	} else {
		roomingHouseID = userPayload.RoomingHouseID
	}

	newRoom := models.Room{
		Name:           roomBody.Name,
		Floor:          roomBody.Floor,
		MaxCapacity:    roomBody.MaxCapacity,
		SizeID:         roomBody.SizeID,
		PackageID:      roomBody.PackageID,
		IsVacant:       true,
		RoomingHouseID: roomingHouseID,
	}

	if err := rc.roomRepo.CreateRoom(&newRoom); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to create room"))
	}

	var roomFacilities []models.RoomFacility
	for _, roomFacilityID := range roomBody.RoomFacilities {
		roomFacility := models.RoomFacility{
			RoomID:     newRoom.ID,
			FacilityID: roomFacilityID,
		}
		roomFacilities = append(roomFacilities, roomFacility)
	}

	if err := rc.roomFacilityRepo.CreateRoomFacility(&roomFacilities); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to create room facility"))
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "success to create room"})
}

// func (rc *RoomController) GetRoomsByRoomingHouseID(c echo.Context) error {
// 	roomingHouseID := c.Param("id")

// 	parsedRoomingHouseID, err := uuid.Parse(roomingHouseID)
// 	if err != nil {
// 		return utils.HandlerError(c, utils.NewBadRequestError("invalid rooming house id"))
// 	}

// 	rooms, err := rc.roomRepo.FindRoomsByRoomingHouseID(parsedRoomingHouseID)
// 	if err != nil {
// 		return utils.HandlerError(c, utils.NewInternalError("failed to get rooms"))
// 	}

// 	return c.JSON(http.StatusOK, rooms)
// }

func (rc *RoomController) GetRoomByID(c echo.Context) error {
	roomID := c.Param("id")

	parsedRoomID, err := uuid.Parse(roomID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid room id"))
	}

	room, err := rc.roomRepo.FindRoomByID(parsedRoomID)
	if err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to get room"))
	}

	return c.JSON(http.StatusOK, room)
}

func (rc *RoomController) GetAllRooms(c echo.Context) error {
	roomingHouseID := c.QueryParam("rooming_house_id")

	if roomingHouseID == "" {
		roomingHouseID = uuid.Nil.String()
	}

	parsedRoomingHouseID, err := uuid.Parse(roomingHouseID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid rooming house id"))
	}

	rooms, err := rc.roomRepo.FindAllRooms(parsedRoomingHouseID)
	if err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to get rooms"))
	}

	return c.JSON(http.StatusOK, rooms)
}

func (rc *RoomController) UpdateRoomByID(c echo.Context) error {
	roomID := c.Param("id")
	var roomBody models.UpdateRoomBody

	if err := c.Bind(&roomBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if roomBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	if roomBody.Floor == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("floor is required"))
	}

	if roomBody.MaxCapacity == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("max capacity is required"))
	}

	if roomBody.SizeID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("size id is required"))
	}

	if roomBody.PackageID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("pricing id is required"))
	}

	if len(roomBody.RoomFacilities) == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("facilities is required"))
	}

	parsedRoomID, err := uuid.Parse(roomID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid room id"))
	}

	if _, err := rc.roomRepo.FindRoomByID(parsedRoomID); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("room not found"))
	}

	updatedRoom := models.Room{
		Name:        roomBody.Name,
		Floor:       roomBody.Floor,
		MaxCapacity: roomBody.MaxCapacity,
		SizeID:      roomBody.SizeID,
		PackageID:   roomBody.PackageID,
	}

	if err := rc.roomRepo.UpdateRoomByID(&updatedRoom, parsedRoomID); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to update room"))
	}

	var roomFacilities []models.RoomFacility
	for _, roomFacilityID := range roomBody.RoomFacilities {
		roomFacility := models.RoomFacility{
			RoomID:     parsedRoomID,
			FacilityID: roomFacilityID,
		}
		roomFacilities = append(roomFacilities, roomFacility)
	}

	if err := rc.roomFacilityRepo.UpdateRoomFacilityByRoomID(&roomFacilities, parsedRoomID); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to update room facility"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "success to update room"})
}

func (rc *RoomController) DeleteRoomByID(c echo.Context) error {
	roomID := c.Param("id")

	parsedRoomID, err := uuid.Parse(roomID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid room id"))
	}

	if _, err := rc.roomRepo.FindRoomByID(parsedRoomID); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("room not found"))
	}

	if err := rc.roomRepo.DeleteRoomByID(parsedRoomID); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to delete room"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "success to delete room"})
}
