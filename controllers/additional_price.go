package controllers

import (
	"net/http"
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AdditionalPriceController struct {
	additionalPriceRepo  repositories.AdditionalPriceRepository
	additionalPeriodRepo repositories.AdditionalPeriodRepository
	periodRepo           repositories.PeriodRepository
}

func NewAdditionalPriceController(additionalPriceRepo repositories.AdditionalPriceRepository, additionalPeriodRepo repositories.AdditionalPeriodRepository, periodRepo repositories.PeriodRepository) *AdditionalPriceController {
	return &AdditionalPriceController{additionalPriceRepo: additionalPriceRepo, additionalPeriodRepo: additionalPeriodRepo, periodRepo: periodRepo}
}

func (apc *AdditionalPriceController) CreateAdditionalPrice(c echo.Context) error {
	var additionalPriceBody models.AddAdditionalPriceBody

	userPayload := c.Get("userPayload").(*models.JWTPayload)

	if err := c.Bind(&additionalPriceBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if additionalPriceBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	if additionalPriceBody.MonthlyPrice == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("monthly price is required"))
	}

	if additionalPriceBody.AnnualPrice == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("annual price is required"))
	}

	var roomingHouseID uuid.UUID

	if userPayload.Role == "owner" {
		roomingHouseID = additionalPriceBody.RoomingHouseID
	} else {
		roomingHouseID = userPayload.RoomingHouseID
	}

	newAdditionalPrice := models.AdditionalPrice{
		Name:           additionalPriceBody.Name,
		RoomingHouseID: roomingHouseID,
	}

	if err := apc.additionalPriceRepo.CreateAdditionalPrice(&newAdditionalPrice); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to create additional price"))
	}

	daily, err := apc.periodRepo.FindPeriodByName("daily")
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("failed to find daily period"))
	}

	weekly, err := apc.periodRepo.FindPeriodByName("weekly")
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("failed to find weekly period"))
	}

	monthly, err := apc.periodRepo.FindPeriodByName("monthly")
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("failed to find monthly period"))
	}

	annual, err := apc.periodRepo.FindPeriodByName("annual")
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("failed to find annual period"))
	}

	additionalPeriods := []models.AdditionalPeriod{
		{
			PeriodID:          daily.ID,
			AdditionalPriceID: newAdditionalPrice.ID,
			Price:             additionalPriceBody.DailyPrice,
		},
		{
			PeriodID:          weekly.ID,
			AdditionalPriceID: newAdditionalPrice.ID,
			Price:             additionalPriceBody.WeeklyPrice,
		},
		{
			PeriodID:          monthly.ID,
			AdditionalPriceID: newAdditionalPrice.ID,
			Price:             additionalPriceBody.MonthlyPrice,
		},
		{
			PeriodID:          annual.ID,
			AdditionalPriceID: newAdditionalPrice.ID,
			Price:             additionalPriceBody.AnnualPrice,
		},
	}

	if err := apc.additionalPeriodRepo.CreateAdditionalPeriod(&additionalPeriods); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to create additional period"))
	}

	return c.JSON(http.StatusCreated, newAdditionalPrice)
}

func (apc *AdditionalPriceController) FindAdditionalPriceByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid id"))
	}

	additionalPrice, err := apc.additionalPriceRepo.FindAdditionalPriceByID(id)
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("additional price not found"))
	}

	return c.JSON(http.StatusOK, additionalPrice)
}

func (apc *AdditionalPriceController) FindAllAdditionalPrices(c echo.Context) error {
	additionalPrices, err := apc.additionalPriceRepo.FindAllAdditionalPrices()
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("additional prices not found"))
	}

	return c.JSON(http.StatusOK, additionalPrices)
}

func (apc *AdditionalPriceController) UpdateAdditionalPriceByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid id"))
	}

	var additionalPriceBody models.UpdateAdditionalPriceBody

	if err := c.Bind(&additionalPriceBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if additionalPriceBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	if additionalPriceBody.MonthlyPrice == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("monthly price is required"))
	}

	if additionalPriceBody.AnnualPrice == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("annual price is required"))
	}

	updatedAdditionalPrice := models.AdditionalPrice{
		Name: additionalPriceBody.Name,
	}

	if err := apc.additionalPriceRepo.UpdateAdditionalPriceByID(&updatedAdditionalPrice, id); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to update additional price"))
	}

	daily, err := apc.periodRepo.FindPeriodByName("daily")
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("failed to find daily period"))
	}

	weekly, err := apc.periodRepo.FindPeriodByName("weekly")
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("failed to find weekly period"))
	}

	monthly, err := apc.periodRepo.FindPeriodByName("monthly")
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("failed to find monthly period"))
	}

	annually, err := apc.periodRepo.FindPeriodByName("annual")
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("failed to find annual period"))
	}

	additionalPeriods := []models.AdditionalPeriod{
		{
			PeriodID:          daily.ID,
			AdditionalPriceID: id,
			Price:             additionalPriceBody.DailyPrice,
		},
		{
			PeriodID:          weekly.ID,
			AdditionalPriceID: id,
			Price:             additionalPriceBody.WeeklyPrice,
		},
		{
			PeriodID:          monthly.ID,
			AdditionalPriceID: id,
			Price:             additionalPriceBody.MonthlyPrice,
		},
		{
			PeriodID:          annually.ID,
			AdditionalPriceID: id,
			Price:             additionalPriceBody.AnnualPrice,
		},
	}

	if err := apc.additionalPeriodRepo.UpdateAdditionalPeriod(&additionalPeriods, id); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to update additional period"))
	}

	return c.JSON(http.StatusOK, updatedAdditionalPrice)
}

func (apc *AdditionalPriceController) DeleteAdditionalPriceByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid id"))
	}

	if err := apc.additionalPriceRepo.DeleteAdditionalPriceByID(id); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to delete additional price"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "additional price deleted"})

}
