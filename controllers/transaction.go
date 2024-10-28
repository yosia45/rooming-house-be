package controllers

import (
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	transactionRepo         repositories.TransactionRepository
	transactionCategoryRepo repositories.TransactionCategoryRepository
	tenantRepo              repositories.TenantRepository
	roomRepo                repositories.RoomRepository
	periodPackageRepo       repositories.PeriodPackageRepository
	periodRepo              repositories.PeriodRepository
}

func NewTransactionController(transactionRepo repositories.TransactionRepository, transactionCategoryRepo repositories.TransactionCategoryRepository, tenantRepo repositories.TenantRepository, periodPackageRepo repositories.PeriodPackageRepository, periodRepo repositories.PeriodRepository, roomRepo repositories.RoomRepository) *TransactionController {
	return &TransactionController{transactionRepo: transactionRepo, transactionCategoryRepo: transactionCategoryRepo, tenantRepo: tenantRepo, periodPackageRepo: periodPackageRepo, periodRepo: periodRepo, roomRepo: roomRepo}
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

	if transactionBody.IsRoom {
		if transactionBody.RoomID == uuid.Nil {
			return utils.HandlerError(c, utils.NewBadRequestError("room id is required"))
		}
	}

	var amount float64

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

		tenant, err := tc.tenantRepo.FindTenantByID(transactionBody.TenantID)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant not found"))
		}

		room, err := tc.roomRepo.FindRoomByID(tenant.RoomID)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("room not found"))
		}

		periodPackage, err := tc.periodPackageRepo.FindPeriodPackageByPeriodIDPackageID(tenant.PeriodID, room.PackageID)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("period package not found"))
		}

		if len(tenant.AdditionalPrices) > 0 {
			for _, additionalPrice := range tenant.AdditionalPrices {
				amount += additionalPrice.AdditionalPeriods[0].Price
			}
		}

		amount += periodPackage.Price * float64(tenant.RegularPaymentDuration)

		if err := tc.transactionRepo.CreateTransaction(&models.Transaction{
			Day:                   transactionBody.Day,
			Month:                 transactionBody.Month,
			Year:                  transactionBody.Year,
			Amount:                amount,
			TransactionCategoryID: transactionBody.TransactionCategoryID,
			TenantID:              tenant.ID,
			RoomID:                tenant.RoomID,
			RoomingHouseID:        tenant.RoomingHouseID,
		}); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to create transaction"))
		}

		var endDate time.Time

		period, err := tc.periodRepo.FindPeriodByID(tenant.PeriodID)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("period not found"))
		}

		if period.Name == "Monthly" {
			endDate = time.Date(transactionBody.Year, time.Month(transactionBody.Month), transactionBody.Day, 0, 0, 0, 0, time.UTC).AddDate(0, tenant.RegularPaymentDuration, 0)
		} else if period.Name == "Annually" {
			endDate = time.Date(transactionBody.Year, time.Month(transactionBody.Month), transactionBody.Day, 0, 0, 0, 0, time.UTC).AddDate(tenant.RegularPaymentDuration, 0, 0)
		} else if period.Name == "Daily" {
			endDate = time.Date(transactionBody.Year, time.Month(transactionBody.Month), transactionBody.Day, 0, 0, 0, 0, time.UTC).AddDate(0, 0, tenant.RegularPaymentDuration)
		} else if period.Name == "Weekly" {
			endDate = time.Date(transactionBody.Year, time.Month(transactionBody.Month), transactionBody.Day, 0, 0, 0, 0, time.UTC).AddDate(0, 0, tenant.RegularPaymentDuration*7)
		}

		if err := tc.tenantRepo.UpdateTenantByID(&models.Tenant{
			StartDate: time.Date(transactionBody.Year, time.Month(transactionBody.Month), transactionBody.Day, 0, 0, 0, 0, time.UTC),
			EndDate:   endDate,
		}, tenant.ID); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to update tenant"))
		}
	} else if transactionCategory.Name == "Deposit" {
		if transactionBody.TenantID == uuid.Nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant id is required"))
		}

		if transactionBody.Amount == 0 {
			return utils.HandlerError(c, utils.NewBadRequestError("amount is required"))
		}

		tenant, err := tc.tenantRepo.FindTenantByID(transactionBody.TenantID)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant not found"))
		}

		if err := tc.transactionRepo.CreateTransaction(&models.Transaction{
			Day:                   transactionBody.Day,
			Month:                 transactionBody.Month,
			Year:                  transactionBody.Year,
			Amount:                transactionBody.Amount,
			TransactionCategoryID: transactionBody.TransactionCategoryID,
			TenantID:              transactionBody.TenantID,
			RoomingHouseID:        tenant.RoomingHouseID,
		}); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to create transaction"))
		}

		if err := tc.tenantRepo.UpdateTenantByID(&models.Tenant{
			IsDepositPaid: true,
		}, tenant.ID); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to update tenant"))
		}
	} else if transactionCategory.Name == "Deposit Payback" {
		if transactionBody.TenantID == uuid.Nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant id is required"))
		}

		if transactionBody.Amount == 0 {
			return utils.HandlerError(c, utils.NewBadRequestError("amount is required"))
		}

		tenant, err := tc.tenantRepo.FindTenantByID(transactionBody.TenantID)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant not found"))
		}

		if err := tc.transactionRepo.CreateTransaction(&models.Transaction{
			Day:                   transactionBody.Day,
			Month:                 transactionBody.Month,
			Year:                  transactionBody.Year,
			Amount:                transactionBody.Amount,
			TransactionCategoryID: transactionBody.TransactionCategoryID,
			TenantID:              transactionBody.TenantID,
			RoomingHouseID:        tenant.RoomingHouseID,
		}); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to create transaction"))
		}

		if err := tc.tenantRepo.UpdateTenantByID(&models.Tenant{
			IsDepositBack: true,
		}, tenant.ID); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to update tenant"))
		}
	} else {
		if transactionBody.RoomingHouseID == uuid.Nil && transactionBody.RoomID == uuid.Nil {
			return utils.HandlerError(c, utils.NewBadRequestError("rooming house id or room id is required"))
		}

		if transactionBody.Amount == 0 {
			return utils.HandlerError(c, utils.NewBadRequestError("amount is required"))
		}

		if err := tc.transactionRepo.CreateTransaction(&models.Transaction{
			Day:                   transactionBody.Day,
			Month:                 transactionBody.Month,
			Year:                  transactionBody.Year,
			Amount:                transactionBody.Amount,
			TransactionCategoryID: transactionBody.TransactionCategoryID,
			RoomID:                transactionBody.RoomID,
			RoomingHouseID:        transactionBody.RoomingHouseID,
		}); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to create transaction"))
		}
	}

	return c.JSON(200, map[string]interface{}{
		"message": "transaction created successfully",
	})
}
