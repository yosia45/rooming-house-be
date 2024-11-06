package controllers

import (
	"net/http"
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SizeController struct {
	sizeRepo         repositories.SizeRepository
	roomingHouseRepo repositories.RoomingHouseRepository
}

func NewSizeController(sizeRepo repositories.SizeRepository, roomingHouseRepo repositories.RoomingHouseRepository) *SizeController {
	return &SizeController{sizeRepo: sizeRepo, roomingHouseRepo: roomingHouseRepo}
}

func (sc *SizeController) CreateSize(c echo.Context) error {
	var sizeBody models.AddSizeBody
	userPayload := c.Get("userPayload").(*models.JWTPayload)

	if err := c.Bind(&sizeBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if sizeBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	if sizeBody.Long == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("long is required"))
	}

	if sizeBody.Width == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("width is required"))
	}

	if userPayload.Role == "owner" && sizeBody.RoomingHouseID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("rooming house id is required"))
	}

	var roomingHouseID uuid.UUID

	if userPayload.Role == "owner" {
		roomingHouseID = sizeBody.RoomingHouseID
	} else {
		roomingHouseID = userPayload.RoomingHouseID
	}

	newSize := models.Size{
		Name:           sizeBody.Name,
		Long:           sizeBody.Long,
		Width:          sizeBody.Width,
		RoomingHouseID: roomingHouseID,
	}

	if err := sc.sizeRepo.CreateSize(&newSize); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to create size"))
	}

	return c.JSON(http.StatusCreated, newSize)
}

func (sc *SizeController) FindSizeByID(c echo.Context) error {
	sizeID := c.Param("id")

	parseSizeID, err := uuid.Parse(sizeID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid size ID"))
	}

	size, err := sc.sizeRepo.FindSizeByID(parseSizeID)
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("size not found"))
	}

	return c.JSON(http.StatusOK, size)
}

func (sc *SizeController) FindAllSizes(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)
	filteredRoomingHouseID := c.QueryParam("roomingHouseID")

	var roomingHouseIDs []uuid.UUID

	if userPayload.Role == "admin" {
		roomingHouseIDs = append(roomingHouseIDs, userPayload.RoomingHouseID)
	} else {
		if filteredRoomingHouseID == "" {
			IDs, err := sc.roomingHouseRepo.FindAllRoomingHouse(userPayload.RoomingHouseID, userPayload.UserID, userPayload.Role)
			if err != nil {
				return utils.HandlerError(c, utils.NewInternalError("failed to get rooming house"))
			}

			for _, roomingHouseID := range IDs {
				roomingHouseIDs = append(roomingHouseIDs, roomingHouseID.ID)
			}
		} else {
			parseRoomingHouseID, err := uuid.Parse(filteredRoomingHouseID)
			if err != nil {
				return utils.HandlerError(c, utils.NewBadRequestError("invalid rooming house ID"))
			}

			roomingHouseIDs = append(roomingHouseIDs, parseRoomingHouseID)
		}
	}

	sizes, err := sc.sizeRepo.FindAllSizes(roomingHouseIDs)
	if err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to get sizes"))
	}

	return c.JSON(http.StatusOK, sizes)
}

func (sc *SizeController) UpdateSizeByID(c echo.Context) error {
	var sizeBody models.UpdateSizeBody
	sizeID := c.Param("id")

	if err := c.Bind(&sizeBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if sizeBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	if sizeBody.Long == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("long is required"))
	}

	if sizeBody.Width == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("width is required"))
	}

	parsedSizeID, err := uuid.Parse(sizeID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid size ID"))
	}

	newSize := models.Size{
		Name:  sizeBody.Name,
		Long:  sizeBody.Long,
		Width: sizeBody.Width,
	}

	if err := sc.sizeRepo.UpdateSizeByID(&newSize, parsedSizeID); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to update size"))
	}

	return c.JSON(http.StatusOK, newSize)
}

func (sc *SizeController) DeleteSizeByID(c echo.Context) error {
	sizeID := c.Param("id")

	parsedSizeID, err := uuid.Parse(sizeID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid size ID"))
	}

	if err := sc.sizeRepo.DeleteSizeByID(parsedSizeID); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to delete size"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "size deleted"})
}
