package controllers

import (
	"net/http"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/labstack/echo/v4"
)

type PeriodController struct {
	periodRepo repositories.PeriodRepository
}

func NewPeriodController(periodRepo repositories.PeriodRepository) *PeriodController {
	return &PeriodController{periodRepo: periodRepo}
}

func (pc *PeriodController) GetAllPeriods(c echo.Context) error {

	periods, err := pc.periodRepo.FindAllPeriods()
	if err != nil {
		return utils.HandlerError(c, utils.NewInternalError(err.Error()))
	}

	return c.JSON(http.StatusOK, periods)
}
