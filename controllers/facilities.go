package controllers

import (
	"net/http"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/labstack/echo/v4"
)

type FacilityController struct {
	facilityRepo repositories.FacilityRepository
}

func NewFacilityController(facilityRepo repositories.FacilityRepository) *FacilityController {
	return &FacilityController{facilityRepo: facilityRepo}
}

func (fc *FacilityController) GetAllFacilities(c echo.Context) error {
	facilities, err := fc.facilityRepo.GetAllFacilities()
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError(err.Error()))
	}

	return c.JSON(http.StatusOK, facilities)
}
