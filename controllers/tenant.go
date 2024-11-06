package controllers

import (
	"fmt"
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

	if tenantBody.IsTenant == nil {
		return utils.HandlerError(c, utils.NewBadRequestError("is tenant is required"))
	}

	if *tenantBody.IsTenant && tenantBody.EmergencyContact == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("emergency contact is required"))
	}

	if tenantBody.RoomID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("room id is required"))
	}

	if tenantBody.PeriodID == uuid.Nil && *tenantBody.IsTenant {
		return utils.HandlerError(c, utils.NewBadRequestError("period id is required"))
	}

	if tenantBody.RegularPaymentDuration == 0 && *tenantBody.IsTenant {
		return utils.HandlerError(c, utils.NewBadRequestError("regular payment duration is required"))
	}

	if !*tenantBody.IsTenant && tenantBody.TenantID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("tenant id is required"))
	}

	if *tenantBody.IsTenant {
		tenantBody.TenantID = uuid.Nil
	}

	var roomingHouseID uuid.UUID

	if userPayload.Role == "owner" {
		roomingHouseID = tenantBody.RoomingHouseID
	} else {
		roomingHouseID = userPayload.RoomingHouseID
	}

	_, err := tc.roomingHouseRepo.FindRoomingHouseByID(roomingHouseID, userPayload.UserID, userPayload.Role)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("rooming house not found"))
	}

	room, err := tc.roomRepo.FindRoomByID(tenantBody.RoomID, roomingHouseID, userPayload.UserID, userPayload.Role)
	if err != nil {
		if err.Error() == "record not found" {
			return utils.HandlerError(c, utils.NewBadRequestError("room not found"))
		}
	}

	if room.RoomingHouseID != roomingHouseID {
		return utils.HandlerError(c, utils.NewBadRequestError("room not found"))
	}

	if !room.IsVacant {
		if !*tenantBody.IsTenant {
			return utils.HandlerError(c, utils.NewBadRequestError("room is not vacant"))
		} else {
			if len(room.Tenants) == room.MaxCapacity {
				return utils.HandlerError(c, utils.NewBadRequestError("room is full"))
			}
		}
	} else {
		if !*tenantBody.IsTenant {
			hasTenant := false
			for _, tenant := range room.Tenants {
				if tenant.IsTenant {
					hasTenant = true
					break
				}
			}

			if !hasTenant {
				return utils.HandlerError(c, utils.NewBadRequestError("need a tenant"))
			}
		}
	}

	newTenant := models.Tenant{
		Name:                   tenantBody.Name,
		Gender:                 tenantBody.Gender,
		PhoneNumber:            tenantBody.PhoneNumber,
		EmergencyContact:       tenantBody.EmergencyContact,
		IsTenant:               *tenantBody.IsTenant,
		RegularPaymentDuration: tenantBody.RegularPaymentDuration,
		RoomingHouseID:         roomingHouseID,
		RoomID:                 tenantBody.RoomID,
		PeriodID:               tenantBody.PeriodID,
		TenantID:               (uuid.UUID)(tenantBody.TenantID),
	}

	if err := tc.tenantRepo.CreateTenant(&newTenant); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to create tenant"))
	}

	fmt.Println(len(tenantBody.TenantAdditionalIDs))

	if len(tenantBody.TenantAdditionalIDs) > 0 {
		var tenantAdditionalPrices []models.TenantAdditionalPrice
		for _, tenantAdditionalID := range tenantBody.TenantAdditionalIDs {
			tenantAdditionalPrice := models.TenantAdditionalPrice{
				TenantID:          newTenant.ID,
				AdditionalPriceID: tenantAdditionalID,
			}

			tenantAdditionalPrices = append(tenantAdditionalPrices, tenantAdditionalPrice)
		}

		fmt.Println(tenantAdditionalPrices)

		if err := tc.tenantAdditionalPriceRepo.CreateTenantAdditional(&tenantAdditionalPrices); err != nil {
			return utils.HandlerError(c, utils.NewBadRequestError("failed to create tenant additional prices"))
		}
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "success to create tenant"})
}

func (tc *TenantController) FindAllTenants(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)

	roomingHouseID := c.QueryParam("roomingHouseID")
	isTenant := c.QueryParam("isTenant")

	var roomingHouseIDs []uuid.UUID

	if userPayload.Role == "admin" {
		roomingHouseIDs = append(roomingHouseIDs, userPayload.RoomingHouseID)
	} else {
		if roomingHouseID != "" {
			parsedRoomingHouseID, err := uuid.Parse(roomingHouseID)
			if err != nil {
				return utils.HandlerError(c, utils.NewBadRequestError("invalid rooming house id"))
			}

			roomingHouseIDs = append(roomingHouseIDs, parsedRoomingHouseID)
		} else {
			roomingHouses, err := tc.roomingHouseRepo.FindAllRoomingHouse(uuid.Nil, userPayload.UserID, userPayload.Role)
			if err != nil {
				return utils.HandlerError(c, utils.NewBadRequestError("failed to find rooming houses"))
			}

			for _, ID := range roomingHouses {
				roomingHouseIDs = append(roomingHouseIDs, ID.ID)
			}

		}
	}

	var isTenantBool bool

	if isTenant == "" || isTenant == "false" {
		isTenantBool = false
	} else {
		isTenantBool = true
	}

	tenants, err := tc.tenantRepo.FindAllTenants(roomingHouseIDs, isTenantBool)
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
