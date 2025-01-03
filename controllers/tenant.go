package controllers

import (
	"net/http"
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TenantController struct {
	tenantRepo                repositories.TenantRepository
	tenantAdditionalPriceRepo repositories.TenantAdditionalRepository
	roomingHouseRepo          repositories.RoomingHouseRepository
	roomRepo                  repositories.RoomRepository
}

func NewTenantController(tenantRepo repositories.TenantRepository, tenantAdditionalRepo repositories.TenantAdditionalRepository, roomingHouseRepo repositories.RoomingHouseRepository, roomRepo repositories.RoomRepository) *TenantController {
	return &TenantController{tenantRepo: tenantRepo, tenantAdditionalPriceRepo: tenantAdditionalRepo, roomingHouseRepo: roomingHouseRepo, roomRepo: roomRepo}
}

func (tc *TenantController) CreateTenant(c echo.Context) error {
	var tenantBody models.AddTenantBody
	userPayload := c.Get("userPayload").(*models.JWTPayload)

	if err := c.Bind(&tenantBody); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if tenantBody.Name == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("name is required"))
	}

	if tenantBody.Gender == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("gender is required"))
	}

	if tenantBody.PhoneNumber == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("phone number is required"))
	}

	var roomingHouseID uuid.UUID

	if userPayload.Role == "owner" {
		roomingHouseID = tenantBody.RoomingHouseID
	} else {
		roomingHouseID = userPayload.RoomingHouseID
	}

	if tenantBody.IsTenant {
		if tenantBody.EmergencyContact == "" {
			return utils.HandlerError(c, utils.NewBadRequestError("emergency contact is required"))
		}

		if tenantBody.RoomID == nil {
			return utils.HandlerError(c, utils.NewBadRequestError("room id is required"))
		}

		if tenantBody.PeriodID == nil {
			return utils.HandlerError(c, utils.NewBadRequestError("period id is required"))
		}

		if tenantBody.RegularPaymentDuration == 0 {
			return utils.HandlerError(c, utils.NewBadRequestError("regular payment duration is required"))
		}

		var room *models.RoomDetailResponse
		var err error

		room, err = tc.roomRepo.FindRoomByID(*tenantBody.RoomID, roomingHouseID, userPayload.UserID, userPayload.Role)
		if err != nil {
			if err.Error() == "record not found" {
				return utils.HandlerError(c, utils.NewBadRequestError("room not found"))
			}
		}

		if room.RoomingHouseID != roomingHouseID {
			return utils.HandlerError(c, utils.NewBadRequestError("room not found"))
		}

		tenantBody.TenantID = uuid.Nil
	} else {
		if tenantBody.TenantID == uuid.Nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant id is required"))
		}

		tenant, err := tc.tenantRepo.FindTenantByID(tenantBody.TenantID, []uuid.UUID{roomingHouseID})
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("tenant not found"))
		}

		room, err := tc.roomRepo.FindRoomByID(tenant.BookedRoomID, roomingHouseID, userPayload.UserID, userPayload.Role)
		if err != nil {
			if err.Error() == "record not found" {
				return utils.HandlerError(c, utils.NewBadRequestError("room not found"))
			}
		}

		if room.MaxCapacity-1 < len(tenant.TenantAssists) {
			return utils.HandlerError(c, utils.NewBadRequestError("room is full"))
		}

		tenantBody.RoomID = nil
		tenantBody.PeriodID = nil
	}

	_, err := tc.roomingHouseRepo.FindRoomingHouseByID(roomingHouseID, userPayload.UserID, userPayload.Role)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("rooming house not found"))
	}

	newTenant := models.Tenant{
		Name:                   tenantBody.Name,
		Gender:                 tenantBody.Gender,
		PhoneNumber:            tenantBody.PhoneNumber,
		EmergencyContact:       tenantBody.EmergencyContact,
		IsTenant:               tenantBody.IsTenant,
		RegularPaymentDuration: tenantBody.RegularPaymentDuration,
		RoomingHouseID:         roomingHouseID,
		RoomID:                 tenantBody.RoomID,
		PeriodID:               tenantBody.PeriodID,
		TenantID:               (uuid.UUID)(tenantBody.TenantID),
	}

	if err := tc.tenantRepo.CreateTenant(&newTenant); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to create tenant"))
	}

	if len(tenantBody.TenantAdditionalIDs) > 0 {
		var tenantAdditionalPrices []models.TenantAdditionalPrice
		for _, tenantAdditionalID := range tenantBody.TenantAdditionalIDs {
			tenantAdditionalPrice := models.TenantAdditionalPrice{
				TenantID:          newTenant.ID,
				AdditionalPriceID: tenantAdditionalID,
			}

			tenantAdditionalPrices = append(tenantAdditionalPrices, tenantAdditionalPrice)
		}

		if err := tc.tenantAdditionalPriceRepo.CreateTenantAdditional(&tenantAdditionalPrices); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to create tenant additional prices"))
		}
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "success to create tenant"})
}

func (tc *TenantController) FindAllTenants(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)
	is_tenant := c.QueryParam("is_tenant")
	var IsTenant bool

	var roomingHouseIDs []uuid.UUID
	if is_tenant == "true" {
		IsTenant = true
	} else {
		IsTenant = false
	}

	if userPayload.Role == "admin" {
		roomingHouseIDs = append(roomingHouseIDs, userPayload.RoomingHouseID)
	} else {
		roomingHouses, err := tc.roomingHouseRepo.FindAllRoomingHouse(uuid.Nil, userPayload.UserID, userPayload.Role)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to find rooming houses"))
		}

		for _, ID := range roomingHouses {
			roomingHouseIDs = append(roomingHouseIDs, ID.ID)
		}
	}

	tenants, err := tc.tenantRepo.FindAllTenants(roomingHouseIDs, IsTenant)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find tenants"))
	}

	return c.JSON(http.StatusOK, tenants)
}

func (tc *TenantController) FindTenantByID(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)

	tenantID := c.Param("id")

	parsedTenantID, err := uuid.Parse(tenantID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid tenant id"))
	}

	var roomingHouseIDs []uuid.UUID

	if userPayload.Role == "admin" {
		roomingHouseIDs = append(roomingHouseIDs, userPayload.RoomingHouseID)
	} else {
		roomingHouses, err := tc.roomingHouseRepo.FindAllRoomingHouse(uuid.Nil, userPayload.UserID, userPayload.Role)
		if err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to find rooming houses"))
		}

		for _, ID := range roomingHouses {
			roomingHouseIDs = append(roomingHouseIDs, ID.ID)
		}

	}

	tenant, err := tc.tenantRepo.FindTenantByID(parsedTenantID, roomingHouseIDs)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find tenant"))
	}

	return c.JSON(http.StatusOK, tenant)
}

func (tc *TenantController) DeleteTenantByID(c echo.Context) error {
	tenantID := c.Param("id")

	parsedTenantID, err := uuid.Parse(tenantID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid tenant id"))
	}

	if err := tc.tenantRepo.DeleteTenantByID(parsedTenantID); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to delete tenant"))
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "success to delete tenant"})
}
