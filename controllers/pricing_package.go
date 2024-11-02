package controllers

import (
	"net/http"
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PricingPackageController struct {
	pricingPackageRepo repositories.PricingPackageRepository
	periodRepo         repositories.PeriodRepository
	periodPackageRepo  repositories.PeriodPackageRepository
	roomingHouseRepo   repositories.RoomingHouseRepository
}

func NewPricingPackageController(pricingPackageRepo repositories.PricingPackageRepository, periodRepo repositories.PeriodRepository, periodPackageRepo repositories.PeriodPackageRepository, roomingHouseRepo repositories.RoomingHouseRepository) *PricingPackageController {
	return &PricingPackageController{pricingPackageRepo: pricingPackageRepo, periodRepo: periodRepo, periodPackageRepo: periodPackageRepo, roomingHouseRepo: roomingHouseRepo}
}

func (ppc *PricingPackageController) CreatePricingPackage(c echo.Context) error {
	var pricingPackageBody models.AddPricingPackageBody
	userPayload := c.Get("userPayload").(*models.JWTPayload)

	if err := c.Bind(&pricingPackageBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if pricingPackageBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	if pricingPackageBody.RoomingHouseID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("rooming house id is required"))
	}

	if pricingPackageBody.DailyPrice == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("daily price is required"))
	}

	if pricingPackageBody.WeeklyPrice == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("weekly price is required"))
	}

	if pricingPackageBody.MonthlyPrice == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("monthly price is required"))
	}

	if pricingPackageBody.AnnualPrice == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("yearly price is required"))
	}

	if userPayload.Role == "owner" && pricingPackageBody.RoomingHouseID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("rooming house id is required"))
	}

	var roomingHouseID uuid.UUID

	if userPayload.Role == "owner" {
		roomingHouseID = pricingPackageBody.RoomingHouseID
	} else {
		roomingHouseID = userPayload.RoomingHouseID
	}

	newPricingPackage := models.PricingPackage{
		Name:           pricingPackageBody.Name,
		RoomingHouseID: roomingHouseID,
	}

	if err := ppc.pricingPackageRepo.CreatePricingPackage(&newPricingPackage); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to create pricing package"))
	}

	daily, err := ppc.periodRepo.FindPeriodByName("Daily")
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find daily id"))
	}

	weekly, err := ppc.periodRepo.FindPeriodByName("Weekly")
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find weekly id"))
	}

	monthly, err := ppc.periodRepo.FindPeriodByName("Monthly")
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find monthly id"))
	}

	annual, err := ppc.periodRepo.FindPeriodByName("Annually")
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find annual id"))
	}

	periodPackage := []models.PeriodPackage{
		{
			PricingPackageID: newPricingPackage.ID,
			PeriodID:         daily.ID,
			Price:            pricingPackageBody.DailyPrice,
		},
		{
			PricingPackageID: newPricingPackage.ID,
			PeriodID:         weekly.ID,
			Price:            pricingPackageBody.WeeklyPrice,
		},
		{
			PricingPackageID: newPricingPackage.ID,
			PeriodID:         monthly.ID,
			Price:            pricingPackageBody.MonthlyPrice,
		},
		{
			PricingPackageID: newPricingPackage.ID,
			PeriodID:         annual.ID,
			Price:            pricingPackageBody.AnnualPrice,
		},
	}

	if err := ppc.periodPackageRepo.CreatePeriodPackage(&periodPackage); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to create period package"))
	}

	return c.JSON(http.StatusCreated, newPricingPackage)
}

func (ppc *PricingPackageController) GetAllPricingPackages(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)

	var roomingHouseIDs []uuid.UUID

	if userPayload.Role == "admin" {
		roomingHouseIDs = append(roomingHouseIDs, userPayload.RoomingHouseID)
	} else {
		IDs, err := ppc.roomingHouseRepo.FindAllRoomingHouse(userPayload.RoomingHouseID, userPayload.UserID, userPayload.Role)

		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to get rooming house"))
		}

		for _, roomingHouseID := range IDs {
			roomingHouseIDs = append(roomingHouseIDs, roomingHouseID.ID)
		}
	}

	pricingPackages, err := ppc.pricingPackageRepo.FindAllPricingPackages(roomingHouseIDs)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to get all pricing packages"))
	}

	return c.JSON(http.StatusOK, pricingPackages)
}

func (ppc *PricingPackageController) UpdatePricingPackage(c echo.Context) error {
	var pricingPackageBody models.UpdatePricingPackageBody
	pricingPackageID := c.Param("id")

	if err := c.Bind(&pricingPackageBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if pricingPackageBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	if pricingPackageBody.DailyPrice == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("daily price is required"))
	}

	if pricingPackageBody.WeeklyPrice == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("weekly price is required"))
	}

	if pricingPackageBody.MonthlyPrice == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("monthly price is required"))
	}

	if pricingPackageBody.AnnualPrice == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("yearly price is required"))
	}

	pricingPackageUUID, err := uuid.Parse(pricingPackageID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid pricing package id"))
	}

	pricingPackage, err := ppc.pricingPackageRepo.FindPricingPackageByID(pricingPackageUUID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("pricing package not found"))
	}

	pricingPackage.Name = pricingPackageBody.Name

	if err := ppc.pricingPackageRepo.UpdatePricingPackageByID(pricingPackage, pricingPackageUUID); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to update pricing package"))
	}

	daily, err := ppc.periodRepo.FindPeriodByName("Daily")
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find daily id"))
	}

	weekly, err := ppc.periodRepo.FindPeriodByName("Weekly")
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find weekly id"))
	}

	monthly, err := ppc.periodRepo.FindPeriodByName("Monthly")
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find monthly id"))
	}

	annual, err := ppc.periodRepo.FindPeriodByName("Annually")
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find annual id"))
	}

	periodPackage := []models.PeriodPackage{
		{
			PricingPackageID: pricingPackage.ID,
			PeriodID:         daily.ID,
			Price:            pricingPackageBody.DailyPrice,
		},
		{
			PricingPackageID: pricingPackage.ID,
			PeriodID:         weekly.ID,
			Price:            pricingPackageBody.WeeklyPrice,
		},
		{
			PricingPackageID: pricingPackage.ID,
			PeriodID:         monthly.ID,
			Price:            pricingPackageBody.MonthlyPrice,
		},
		{
			PricingPackageID: pricingPackage.ID,
			PeriodID:         annual.ID,
			Price:            pricingPackageBody.AnnualPrice,
		},
	}

	if err := ppc.periodPackageRepo.UpdatePeriodPackageByPackageID(periodPackage, pricingPackageUUID); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to update period package"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "pricing package updated"})
}

func (ppc *PricingPackageController) DeletePricingPackage(c echo.Context) error {
	pricingPackageID := c.Param("id")
	pricingPackageUUID, err := uuid.Parse(pricingPackageID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid pricing package id"))
	}

	if err := ppc.pricingPackageRepo.DeletePricingPackageByID(pricingPackageUUID); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to delete pricing package"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "pricing package deleted"})
}
