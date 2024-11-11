package controllers

import (
	"fmt"
	"net/http"
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RoomController struct {
	roomRepo         repositories.RoomRepository
	roomFacilityRepo repositories.RoomFacilityRepository
	roomingHouseRepo repositories.RoomingHouseRepository
	sizeRepo         repositories.SizeRepository
	packageRepo      repositories.PricingPackageRepository
	facilityRepo     repositories.FacilityRepository
}

func NewRoomController(roomRepo repositories.RoomRepository, roomFacilityRepo repositories.RoomFacilityRepository, roomingHouseRepo repositories.RoomingHouseRepository, sizeRepo repositories.SizeRepository, packageRepo repositories.PricingPackageRepository, facilityRepo repositories.FacilityRepository) *RoomController {
	return &RoomController{roomRepo: roomRepo, roomFacilityRepo: roomFacilityRepo, roomingHouseRepo: roomingHouseRepo, sizeRepo: sizeRepo, facilityRepo: facilityRepo, packageRepo: packageRepo}
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

	if roomBody.RoomingHouseID == uuid.Nil && userPayload.Role == "owner" {
		return utils.HandlerError(c, utils.NewBadRequestError("rooming house id is required"))
	}

	var roomingHouseID uuid.UUID

	if userPayload.Role == "owner" {
		roomingHouseID = roomBody.RoomingHouseID
	} else {
		roomingHouseID = userPayload.RoomingHouseID
	}

	roomingHouse, err := rc.roomingHouseRepo.FindRoomingHouseByID(roomBody.RoomingHouseID, userPayload.UserID, userPayload.Role)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandlerError(c, utils.NewBadRequestError("rooming house not found"))
		}
		return utils.HandlerError(c, utils.NewInternalError("failed to get rooming house"))
	}

	if roomBody.SizeID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("size id is required"))
	}

	size, err := rc.sizeRepo.FindSizeByID(roomBody.SizeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandlerError(c, utils.NewBadRequestError("size not found"))
		}
		return utils.HandlerError(c, utils.NewInternalError("failed to get size"))
	}

	if size.RoomingHouseID != roomingHouseID {
		return utils.HandlerError(c, utils.NewBadRequestError("size not from this rooming house"))
	}

	if roomBody.PackageID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("package id is required"))
	}

	packagePricing, err := rc.packageRepo.FindPricingPackageByID(roomBody.PackageID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandlerError(c, utils.NewBadRequestError("pricing package not found"))
		}
		return utils.HandlerError(c, utils.NewInternalError("failed to get pricing package"))
	}

	if packagePricing.RoomingHouseID != roomingHouseID {
		return utils.HandlerError(c, utils.NewBadRequestError("pricing package not from this rooming house"))
	}

	if len(roomBody.RoomFacilities) == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("facilities is required"))
	}

	if roomBody.Floor > roomingHouse.FloorTotal {
		return utils.HandlerError(c, utils.NewBadRequestError("floor is greater than floor total"))
	}

	for _, roomFacilityID := range roomBody.RoomFacilities {
		facility, err := rc.facilityRepo.GetFacilityByID(roomFacilityID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.HandlerError(c, utils.NewBadRequestError("facility not found"))
			}
			return utils.HandlerError(c, utils.NewInternalError("failed to get facility"))
		}

		if !facility.IsRoom {
			return utils.HandlerError(c, utils.NewBadRequestError("facility is not room facility"))
		}
	}

	newRoom := models.Room{
		Name:           roomBody.Name,
		Floor:          roomBody.Floor,
		MaxCapacity:    roomBody.MaxCapacity,
		SizeID:         roomBody.SizeID,
		PackageID:      roomBody.PackageID,
		TenantID:       uuid.Nil,
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

func (rc *RoomController) GetRoomByID(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)
	roomID := c.Param("id")

	parsedRoomID, err := uuid.Parse(roomID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid room id"))
	}

	var roomingHouseID uuid.UUID

	if userPayload.Role == "admin" {
		roomingHouseID = userPayload.RoomingHouseID
	} else {
		roomingHouseID = uuid.Nil
	}

	room, err := rc.roomRepo.FindRoomByID(parsedRoomID, roomingHouseID, userPayload.UserID, userPayload.Role)
	if err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to get room"))
	}

	return c.JSON(http.StatusOK, room)
}

func (rc *RoomController) GetAllRooms(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)

	filteredRoomingHouseID := c.QueryParam("rooming_house_id")

	var roomingHouseIDs []uuid.UUID

	if userPayload.Role == "admin" {
		roomingHouseIDs = append(roomingHouseIDs, userPayload.RoomingHouseID)
	} else {
		if filteredRoomingHouseID == "" {
			roomingHouses, err := rc.roomingHouseRepo.FindAllRoomingHouse(userPayload.RoomingHouseID, userPayload.UserID, userPayload.Role)
			if err != nil {
				return utils.HandlerError(c, utils.NewBadRequestError("failed to get rooming house"))
			}

			for _, roomingHouseID := range roomingHouses {
				roomingHouseIDs = append(roomingHouseIDs, roomingHouseID.ID)
			}
		} else {
			parsedRoomingHouseID, err := uuid.Parse(filteredRoomingHouseID)
			if err != nil {
				return utils.HandlerError(c, utils.NewBadRequestError("invalid rooming house id"))
			}

			roomingHouseIDs = append(roomingHouseIDs, parsedRoomingHouseID)
		}
	}

	rooms, err := rc.roomRepo.FindAllRooms(roomingHouseIDs)
	if err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to get rooms"))
	}

	return c.JSON(http.StatusOK, rooms)
}

func (rc *RoomController) UpdateRoomByID(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)
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

	var roomingHouseID uuid.UUID

	if userPayload.Role == "admin" {
		roomingHouseID = userPayload.RoomingHouseID
	} else {
		roomingHouseID = uuid.Nil
	}

	roomByID, err := rc.roomRepo.FindRoomByID(parsedRoomID, roomingHouseID, userPayload.UserID, userPayload.Role)
	if err != nil {
		return utils.HandlerError(c, utils.NewInternalError("room not found"))
	}

	if len(roomByID.Tenants.TenantAssists) > 0 && roomByID.MaxCapacity-1 > roomBody.MaxCapacity {
		return utils.HandlerError(c, utils.NewBadRequestError(fmt.Sprintf("max capacity is less than current tenant count (%d)", len(roomByID.Tenants.TenantAssists)+1)))
	}

	roomingHouse, err := rc.roomingHouseRepo.FindRoomingHouseByID(roomByID.RoomingHouseID, userPayload.UserID, userPayload.Role)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandlerError(c, utils.NewBadRequestError("rooming house not found"))
		}

		return utils.HandlerError(c, utils.NewInternalError("failed to get rooming house"))
	}

	if roomingHouse.FloorTotal < roomBody.Floor {
		return utils.HandlerError(c, utils.NewBadRequestError("floor is greater than floor total"))
	}

	size, err := rc.sizeRepo.FindSizeByID(roomBody.SizeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandlerError(c, utils.NewBadRequestError("size not found"))
		}
		return utils.HandlerError(c, utils.NewInternalError("failed to get size"))
	}

	if size.RoomingHouseID != roomByID.RoomingHouseID {
		return utils.HandlerError(c, utils.NewBadRequestError("size not from this rooming house"))
	}

	packagePricing, err := rc.packageRepo.FindPricingPackageByID(roomBody.PackageID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandlerError(c, utils.NewBadRequestError("pricing package not found"))
		}
		return utils.HandlerError(c, utils.NewInternalError("failed to get pricing package"))
	}

	if packagePricing.RoomingHouseID != roomByID.RoomingHouseID {
		return utils.HandlerError(c, utils.NewBadRequestError("pricing package not from this rooming house"))
	}

	for _, roomFacilityID := range roomBody.RoomFacilities {
		facility, err := rc.facilityRepo.GetFacilityByID(roomFacilityID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.HandlerError(c, utils.NewBadRequestError("facility not found"))
			}
			return utils.HandlerError(c, utils.NewInternalError("failed to get facility"))
		}

		if !facility.IsRoom {
			return utils.HandlerError(c, utils.NewBadRequestError("facility is not room facility"))
		}
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
	userPayload := c.Get("userPayload").(*models.JWTPayload)
	roomID := c.Param("id")

	parsedRoomID, err := uuid.Parse(roomID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid room id"))
	}

	var roomingHouseID uuid.UUID

	if userPayload.Role == "admin" {
		roomingHouseID = userPayload.RoomingHouseID
	} else {
		roomingHouseID = uuid.Nil
	}

	if _, err := rc.roomRepo.FindRoomByID(parsedRoomID, roomingHouseID, userPayload.UserID, userPayload.Role); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("room not found"))
	}

	if err := rc.roomRepo.DeleteRoomByID(parsedRoomID); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to delete room"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "success to delete room"})
}
