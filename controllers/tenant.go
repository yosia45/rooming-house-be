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
	tenantRepo repositories.TenantRepository
}

func NewTenantController(tenantRepo repositories.TenantRepository) *TenantController {
	return &TenantController{tenantRepo: tenantRepo}
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

	if len(tenantBody.TenantAdditionalIDs) == 0 {
		return utils.HandlerError(c, utils.NewBadRequestError("tenant additional ids is required"))
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

	return c.JSON(http.StatusCreated, newTenant)
}
