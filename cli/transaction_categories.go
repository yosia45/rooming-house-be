package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func TransactionCategoryRoutes(e *echo.Echo) {
	transactionCategoryRepo := repositories.NewTransactionCategoryRepository(config.DB)

	transactionCategoryController := controllers.NewTransactionCategoryController(transactionCategoryRepo)

	transactionCategory := e.Group("/transaction-categories")
	transactionCategory.POST("", transactionCategoryController.CreateTransactionCategory)
	transactionCategory.GET("/:id", transactionCategoryController.FindTransactionCategoryByID)
	transactionCategory.GET("", transactionCategoryController.FindAllTransactionCategories)
	transactionCategory.PUT("/:id", transactionCategoryController.UpdateTransactionCategoryByID)
	transactionCategory.DELETE("/:id", transactionCategoryController.DeleteTransactionCategoryByID)
}
