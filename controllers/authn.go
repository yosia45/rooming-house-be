package controllers

import (
	"net/http"
	"rooming-house-cms-be/middlewares"
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/repositories"
	"rooming-house-cms-be/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	ownerRepo        repositories.OwnerRepository
	adminRepo        repositories.AdminRepository
	roomingHouseRepo repositories.RoomingHouseRepository
}

func NewUserController(ownerRepo repositories.OwnerRepository, adminRepo repositories.AdminRepository, roomingHouseRepo repositories.RoomingHouseRepository) *UserController {
	return &UserController{ownerRepo: ownerRepo, adminRepo: adminRepo, roomingHouseRepo: roomingHouseRepo}
}

func (uc *UserController) RegisterOwner(c echo.Context) error {
	var owner models.OwnerRegisterBody
	if err := c.Bind(&owner); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if owner.FullName == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("full name is required"))
	}

	if owner.Email == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("email is required"))
	}

	if owner.Username == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("username is required"))
	}

	if owner.Password == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("password is required"))
	}

	newOwner := models.Owner{
		FullName: owner.FullName,
		Email:    owner.Email,
		Username: owner.Username,
		Password: owner.Password,
		Role:     "owner",
	}

	// Create new user
	if err := uc.ownerRepo.CreateOwner(&newOwner); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to create owner"))
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "success to create owner"})
}

func (uc *UserController) RegisterAdmin(c echo.Context) error {
	userPayload := c.Get("userPayload").(*models.JWTPayload)
	var admin models.AdminRegisterBody
	if err := c.Bind(&admin); err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("invalid input"))
	}

	if admin.FullName == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("full name is required"))
	}

	if admin.Email == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("email is required"))
	}

	if admin.Username == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("username is required"))
	}

	if admin.Password == "" {
		return utils.HandlerError(c, utils.NewBadRequestError("password is required"))
	}

	if admin.RoomingHouseID == uuid.Nil {
		return utils.HandlerError(c, utils.NewBadRequestError("rooming house ID is required"))
	}

	roomingHouse, err := uc.roomingHouseRepo.FindRoomingHouseByID(admin.RoomingHouseID, userPayload.UserID, userPayload.Role)
	if err != nil {
		return utils.HandlerError(c, utils.NewBadRequestError("rooming house not found"))
	}

	if roomingHouse.ID != userPayload.RoomingHouseID {
		return utils.HandlerError(c, utils.NewUnauthorizedError("you are not the owner of this rooming house"))
	}

	newAdmin := models.Admin{
		FullName:       admin.FullName,
		Email:          admin.Email,
		Username:       admin.Username,
		Password:       admin.Password,
		Role:           "admin",
		RoomingHouseID: admin.RoomingHouseID,
	}

	// Create new user

	if err := uc.adminRepo.CreateAdmin(&newAdmin); err != nil {
		return utils.HandlerError(c, utils.NewInternalError("failed to create admin"))
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "success to create admin"})
}

func (uc *UserController) Login(c echo.Context) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	var passwordHash string
	var userID uuid.UUID
	var roomingHouseID uuid.UUID
	var err error

	switch input.Role {
	case "admin":
		admin, err := uc.adminRepo.FindAdminByEmail(input.Email)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "admin not found"})
		}
		passwordHash = admin.Password
		userID = admin.ID
		roomingHouseID = admin.RoomingHouseID
	case "owner":
		owner, err := uc.ownerRepo.FindOwnerByEmail(input.Email)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "owner not found"})
		}
		passwordHash = owner.Password
		userID = owner.ID
		if len(owner.RoomingHouses) > 0 {
			roomingHouseID = owner.RoomingHouses[0].ID
		} else {
			roomingHouseID = uuid.Nil
		}
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid role"})
	}

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}
	// Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid password"})
	}

	// Generate JWT
	token, err := middlewares.GenerateJWT(userID, input.Role, roomingHouseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
