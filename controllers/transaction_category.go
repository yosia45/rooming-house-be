package controllers

import (
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransactionCategoryController struct {
	transactionCategoryRepo repositories.TransactionCategoryRepository
}

func NewTransactionCategoryController(transactionCategoryRepo repositories.TransactionCategoryRepository) *TransactionCategoryController {
	return &TransactionCategoryController{transactionCategoryRepo: transactionCategoryRepo}
}

func (tcc *TransactionCategoryController) CreateTransactionCategory(c echo.Context) error {
	var transactionCategoryBody models.TransactionCategoryBody

	if err := c.Bind(&transactionCategoryBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if transactionCategoryBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	newTransactionCategory := models.TransactionCategory{
		Name:      transactionCategoryBody.Name,
		IsExpense: transactionCategoryBody.IsExpense,
	}

	if err := tcc.transactionCategoryRepo.CreateTransactionCategory(&newTransactionCategory); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to create transaction category"))
	}

	return c.JSON(201, newTransactionCategory)
}

func (tcc *TransactionCategoryController) FindTransactionCategoryByID(c echo.Context) error {
	id := c.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid id"))
	}

	transactionCategory, err := tcc.transactionCategoryRepo.FindTransactionCategoryByID(parsedID)
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("transaction category not found"))
	}

	return c.JSON(200, transactionCategory)
}

func (tcc *TransactionCategoryController) FindAllTransactionCategories(c echo.Context) error {
	transactionCategories, err := tcc.transactionCategoryRepo.FindAllTransactionCategories()
	if err != nil {
		return utils.HandlerError(c, utils.NewNotFoundError("transaction categories not found"))
	}

	return c.JSON(200, transactionCategories)
}

func (tcc *TransactionCategoryController) UpdateTransactionCategoryByID(c echo.Context) error {
	var transactionCategoryBody models.TransactionCategoryBody
	id := c.Param("id")

	if err := c.Bind(&transactionCategoryBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if transactionCategoryBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid id"))
	}

	newTransactionCategory := models.TransactionCategory{
		Name:      transactionCategoryBody.Name,
		IsExpense: transactionCategoryBody.IsExpense,
	}

	if err := tcc.transactionCategoryRepo.UpdateTransactionCategoryByID(&newTransactionCategory, parsedID); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to update transaction category"))
	}

	return c.JSON(200, newTransactionCategory)
}

func (tcc *TransactionCategoryController) DeleteTransactionCategoryByID(c echo.Context) error {
	id := c.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid id"))
	}

	if err := tcc.transactionCategoryRepo.DeleteTransactionCategoryByID(parsedID); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to delete transaction category"))
	}

	return c.NoContent(204)
}
