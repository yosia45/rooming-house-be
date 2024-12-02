package controllers

import (
	"net/http"
	"rooming-house-cms-be/constants"
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"
	"strconv"
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
	roomingHouseRepo        repositories.RoomingHouseRepository
}

func NewTransactionController(transactionRepo repositories.TransactionRepository, transactionCategoryRepo repositories.TransactionCategoryRepository, tenantRepo repositories.TenantRepository, periodPackageRepo repositories.PeriodPackageRepository, periodRepo repositories.PeriodRepository, roomRepo repositories.RoomRepository, roomingHouseRepo repositories.RoomingHouseRepository) *TransactionController {
	return &TransactionController{transactionRepo: transactionRepo, transactionCategoryRepo: transactionCategoryRepo, tenantRepo: tenantRepo, periodPackageRepo: periodPackageRepo, periodRepo: periodRepo, roomRepo: roomRepo, roomingHouseRepo: roomingHouseRepo}
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

	var amount float64

	transactionCategory, err := tc.transactionCategoryRepo.FindTransactionCategoryByID(transactionBody.TransactionCategoryID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("transaction category not found"))
	}

	if transactionCategory.Name == "Rent" {
		if transactionBody.TenantID == nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant id is required"))
		}

		var roomingHouseIDs []uuid.UUID

		roomingHouseIDs = append(roomingHouseIDs, transactionBody.RoomingHouseID)

		tenant, err := tc.tenantRepo.FindTenantByID(*transactionBody.TenantID, roomingHouseIDs)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant not found"))
		}

		room, err := tc.roomRepo.FindRoomByID(tenant.BookedRoomID, tenant.RoomingHouse.ID, userPayload.UserID, userPayload.Role)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("room not found"))
		}

		periodPackage, err := tc.periodPackageRepo.FindPeriodPackageByPeriodIDPackageID(tenant.Period.ID, room.PricingPackage.ID)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("period package not found"))
		}

		if len(tenant.AdditionalPrices) > 0 {
			for _, additionalPrice := range tenant.AdditionalPrices {
				amount += additionalPrice.Price
			}
		}

		amount += periodPackage.Price * float64(tenant.RegularPaymentDuration)

		if err := tc.transactionRepo.CreateTransaction(&models.Transaction{
			Day:                   transactionBody.Day,
			Month:                 transactionBody.Month,
			Year:                  transactionBody.Year,
			Amount:                amount,
			TransactionCategoryID: transactionBody.TransactionCategoryID,
			TenantID:              &tenant.ID,
			RoomID:                &tenant.BookedRoomID,
			RoomingHouseID:        tenant.RoomingHouse.ID,
		}); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to create transaction"))
		}

		var endDate time.Time

		period, err := tc.periodRepo.FindPeriodByID(tenant.Period.ID)
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

		startDate := time.Date(transactionBody.Year, time.Month(transactionBody.Month), transactionBody.Day, 0, 0, 0, 0, time.UTC)

		if err := tc.tenantRepo.UpdateTenantByID(&models.Tenant{
			StartDate: &startDate,
			EndDate:   &endDate,
		}, tenant.ID); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to update tenant"))
		}
	} else if transactionCategory.Name == "Deposit" {
		if transactionBody.TenantID == nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant id is required"))
		}

		if transactionBody.Amount == 0 {
			return utils.HandlerError(c, utils.NewBadRequestError("amount is required"))
		}

		var roomingHouseIDs []uuid.UUID

		roomingHouseIDs = append(roomingHouseIDs, transactionBody.RoomingHouseID)

		tenant, err := tc.tenantRepo.FindTenantByID(*transactionBody.TenantID, roomingHouseIDs)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant not found"))
		}

		if tenant.IsDepositPaid {
			return utils.HandlerError(c, utils.NewBadRequestError("deposit already paid"))
		}

		if err := tc.transactionRepo.CreateTransaction(&models.Transaction{
			Day:                   transactionBody.Day,
			Month:                 transactionBody.Month,
			Year:                  transactionBody.Year,
			Amount:                transactionBody.Amount,
			TransactionCategoryID: transactionBody.TransactionCategoryID,
			RoomID:                &tenant.BookedRoomID,
			TenantID:              transactionBody.TenantID,
			RoomingHouseID:        tenant.RoomingHouse.ID,
		}); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to create transaction"))
		}

		if err := tc.tenantRepo.UpdateTenantByID(&models.Tenant{
			IsDepositPaid: true,
		}, tenant.ID); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to update tenant"))
		}
	} else if transactionCategory.Name == "Deposit Payback" {
		if transactionBody.TenantID == nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant id is required"))
		}

		if transactionBody.Amount == 0 {
			return utils.HandlerError(c, utils.NewBadRequestError("amount is required"))
		}

		var roomingHouseIDs []uuid.UUID

		roomingHouseIDs = append(roomingHouseIDs, transactionBody.RoomingHouseID)

		tenant, err := tc.tenantRepo.FindTenantByID(*transactionBody.TenantID, roomingHouseIDs)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant not found"))
		}

		if !tenant.IsDepositPaid {
			return utils.HandlerError(c, utils.NewBadRequestError("deposit not paid"))
		}

		if err := tc.transactionRepo.CreateTransaction(&models.Transaction{
			Day:                   transactionBody.Day,
			Month:                 transactionBody.Month,
			Year:                  transactionBody.Year,
			Amount:                transactionBody.Amount,
			TransactionCategoryID: transactionBody.TransactionCategoryID,
			TenantID:              transactionBody.TenantID,
			RoomingHouseID:        tenant.RoomingHouse.ID,
		}); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to create transaction"))
		}

		if err := tc.tenantRepo.UpdateTenantByID(&models.Tenant{IsDepositPaid: false, IsDepositBack: true}, tenant.ID); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to update tenant"))
		}
	} else {
		if transactionBody.IsRoom {
			if transactionBody.RoomID == nil {
				return utils.HandlerError(c, utils.NewBadRequestError("room id is required"))
			}
		} else {
			if transactionBody.RoomingHouseID == uuid.Nil {
				return utils.HandlerError(c, utils.NewBadRequestError("rooming house id is required"))
			}
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

func (tc *TransactionController) FindAllTransactions(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)

	// roomingHouseID := c.QueryParam("roomingHouseID")
	var roomingHouseIDs []uuid.UUID

	if userPayload.Role == "admin" {
		roomingHouseIDs = append(roomingHouseIDs, userPayload.RoomingHouseID)
	} else {
		roomingHouses, err := tc.roomingHouseRepo.FindAllRoomingHouse(userPayload.RoomingHouseID, userPayload.UserID, userPayload.Role)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to find rooming houses"))
		}

		for _, roomingHouse := range roomingHouses {
			roomingHouseIDs = append(roomingHouseIDs, roomingHouse.ID)
		}
	}

	transactions, err := tc.transactionRepo.FindAllTransactions(roomingHouseIDs, 0)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find transactions"))
	}

	return c.JSON(http.StatusOK, transactions)
}

func (tc *TransactionController) Dashboard(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)
	roomingHouseID := c.QueryParam("roomingHouseID")
	year := c.QueryParam("year")

	var roomingHouseIDs []uuid.UUID

	now := time.Now()
	currentYear := now.Format("2006")

	if roomingHouseID != "" {
		parsedRoomingHouseID, err := uuid.Parse(roomingHouseID)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("invalid rooming house id"))
		}
		roomingHouseIDs = append(roomingHouseIDs, parsedRoomingHouseID)
	} else {
		roomingHouses, err := tc.roomingHouseRepo.FindAllRoomingHouse(userPayload.RoomingHouseID, userPayload.UserID, userPayload.Role)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to find rooming houses"))
		}

		roomingHouseIDs = append(roomingHouseIDs, roomingHouses[0].ID)
	}

	var yearInt int

	if year == "" {
		parsedYear, err := strconv.Atoi(currentYear)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("invalid year"))
		}

		yearInt = parsedYear
	} else {
		parsedYear, err := strconv.Atoi(year)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("invalid year"))
		}

		yearInt = parsedYear
	}

	transactions, err := tc.transactionRepo.FindAllTransactions(roomingHouseIDs, yearInt)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find transactions"))
	}

	// Group by Rooming House
	groupedData := make(map[string]*models.DashboardData)

	// Iterate through transactions to group data and calculate monthly income/expense
	for _, txn := range *transactions {
		rhName := txn.RoomingHouse.Name

		// Initialize DashboardData if not exists
		if _, exists := groupedData[rhName]; !exists {
			transactionData := make([]models.TransactionDashboardResponse, 12)
			for i := 0; i < 12; i++ {
				transactionData[i] = models.TransactionDashboardResponse{
					Month:   constants.Months[i],
					Year:    txn.Year, // Default year from first transaction
					Index:   i + 1,
					Income:  0,
					Expense: 0,
				}
			}
			groupedData[rhName] = &models.DashboardData{
				RoomingHouseName: rhName,
				TransactionData:  transactionData,
			}
		}

		// Find and update the month in TransactionData
		for i := range groupedData[rhName].TransactionData {
			if groupedData[rhName].TransactionData[i].Index == txn.Month {
				if txn.Category.IsExpense {
					groupedData[rhName].TransactionData[i].Expense += txn.Amount
				} else {
					groupedData[rhName].TransactionData[i].Income += txn.Amount
				}
				break
			}
		}
	}

	// Convert grouped data into final response slice
	finalResponse := []models.DashboardData{}
	for _, data := range groupedData {
		finalResponse = append(finalResponse, *data)
	}

	return c.JSON(http.StatusOK, finalResponse)
}
