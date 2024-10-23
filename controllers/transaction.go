package controllers

import (
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	transactionRepo         repositories.TransactionRepository
	transactionCategoryRepo repositories.TransactionCategoryRepository
}

func NewTransactionController(transactionRepo repositories.TransactionRepository, transactionCategoryRepo repositories.TransactionCategoryRepository) *TransactionController {
	return &TransactionController{transactionRepo: transactionRepo, transactionCategoryRepo: transactionCategoryRepo}
}

func (tc *TransactionController) CreateTransaction(c echo.Context) error {
	var transactionBody models.AddTransactionBody
	userPayload := c.Get("userPayload").(*models.JWTPayload)

	if err := c.Bind(&transactionBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if transactionBody.Day == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("day is required"))
	}

	if transactionBody.Month == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("month is required"))
	}

	if transactionBody.Year == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("year is required"))
	}

	if transactionBody.TransactionCategoryID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("transaction category id is required"))
	}

	if userPayload.Role == "owner" {
		if transactionBody.RoomingHouseID == uuid.Nil {
			return utils.HandlerError(c, utils.NewBadRequestError("rooming house id is required"))
		}
	}

	if transactionBody.IsRoom == true {
		if transactionBody.RoomID == uuid.Nil {
			return utils.HandlerError(c, utils.NewBadRequestError("room id is required"))
		}
	}

	transactionCategory, err := tc.transactionCategoryRepo.FindTransactionCategoryByID(transactionBody.TransactionCategoryID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("transaction category not found"))
	}

	if transactionCategory.Name != "Rent" {
		if transactionBody.TenantID == uuid.Nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant id is required"))
		}

		if transactionBody.Amount == 0 {
			return utils.HandlerError(c, utils.NewBadRequestError("amount is required"))
		}

	}

	return nil
}
