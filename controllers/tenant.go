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
}

func NewTenantController(tenantRepo repositories.TenantRepository, tenantAdditionalRepo repositories.TenantAdditionalRepository) *TenantController {
	return &TenantController{tenantRepo: tenantRepo, tenantAdditionalPriceRepo: tenantAdditionalRepo}
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

	if tenantBody.EmergencyContact == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("emergency contact is required"))
	}

	if tenantBody.RoomID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("room id is required"))
	}

	if tenantBody.PeriodID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("period id is required"))
	}

	if tenantBody.IsTenant == false && tenantBody.TenantID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("tenant id is required"))
	}

	if tenantBody.IsTenant == true {
		tenantBody.TenantID = uuid.Nil
	}

	var roomingHouseID uuid.UUID

	if userPayload.Role == "owner" {
		roomingHouseID = tenantBody.RoomingHouseID
	} else {
		roomingHouseID = userPayload.RoomingHouseID
	}

	newTenant := models.Tenant{
		Name:             tenantBody.Name,
		Gender:           tenantBody.Gender,
		PhoneNumber:      tenantBody.PhoneNumber,
		EmergencyContact: tenantBody.EmergencyContact,
		IsTenant:         tenantBody.IsTenant,
		RoomingHouseID:   roomingHouseID,
		RoomID:           tenantBody.RoomID,
		PeriodID:         tenantBody.PeriodID,
		TenantID:         (uuid.UUID)(tenantBody.TenantID),
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

	return c.JSON(http.StatusCreated, newTenant)
}

func (tc *TenantController) FindAllTenants(c echo.Context) error {
	tenants, err := tc.tenantRepo.FindAllTenants()
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("failed to find tenants"))
	}

	return c.JSON(http.StatusOK, tenants)
}

func (tc *TenantController) FindTenantByID(c echo.Context) error {
	tenantID := c.Param("id")

	parsedTenantID, err := uuid.Parse(tenantID)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid tenant id"))
	}

	tenant, err := tc.tenantRepo.FindTenantByID(parsedTenantID)
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
