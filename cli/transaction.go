package cli

import (
	"rooming-house-cms-be/config"
	"rooming-house-cms-be/controllers"
	"rooming-house-cms-be/middlewares"
	"rooming-house-cms-be/repositories"

	"github.com/labstack/echo/v4"
)

func TransactionRoutes(e *echo.Echo) {
	transactionRepo := repositories.NewTransactionRepository(config.DB)
	transactionCategoryRepo := repositories.NewTransactionCategoryRepository(config.DB)
	tenantRepo := repositories.NewTenantRepository(config.DB)
	roomRepo := repositories.NewRoomRepository(config.DB)
	periodPackageRepo := repositories.NewPeriodPackageRepository(config.DB)
	periodRepo := repositories.NewPeriodRepository(config.DB)
	roomingHouseRepo := repositories.NewRoomingHouseRepository(config.DB)

	transactionController := controllers.NewTransactionController(transactionRepo, transactionCategoryRepo, tenantRepo, periodPackageRepo, periodRepo, roomRepo, roomingHouseRepo)

	transaction := e.Group("/transactions")
	transaction.POST("", transactionController.CreateTransaction, middlewares.JWTAuth)
	transaction.GET("", transactionController.FindAllTransactions, middlewares.JWTAuth)
	transaction.GET("/dashboard", transactionController.Dashboard, middlewares.JWTAuth)
}
