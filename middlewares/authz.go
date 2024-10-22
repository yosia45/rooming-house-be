package middlewares

import (
	"fmt"
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/utils"

	"github.com/labstack/echo/v4"
)

func Authz(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Authz middleware")
		userPayload := c.Get("userPayload").(*models.JWTPayload)

		if userPayload.Role != "owner" {
			return utils.HandlerError(c, utils.NewForbiddenError("only owner can access this resource"))
		}
		return next(c)
	}
}
