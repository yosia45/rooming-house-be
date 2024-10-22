package middlewares

import (
	"fmt"
	"os"
	"rooming-house-cms-be/models"
	"rooming-house-cms-be/utils"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateJWT(userID uuid.UUID, role string, rooming_house_id uuid.UUID) (string, error) {
	var claims jwt.Claims
	if role == "owner" {
		claims = jwt.MapClaims{
			"user_id":          userID,
			"role":             role,
			"rooming_house_id": uuid.Nil,
		}
	} else {
		claims = jwt.MapClaims{
			"user_id":          userID,
			"role":             role,
			"rooming_house_id": rooming_house_id,
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("masuk jwt auth")
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return utils.HandlerError(c, utils.NewUnauthorizedError("please login first"))
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return utils.HandlerError(c, utils.NewUnauthorizedError("invalid authorization header format"))
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			return utils.HandlerError(c, utils.NewUnauthorizedError("invalid token"))
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				return utils.HandlerError(c, utils.NewUnauthorizedError("invalid token: user_id not found"))
			}

			roomingHouseIDStr, ok := claims["rooming_house_id"].(string)
			if !ok {
				return utils.HandlerError(c, utils.NewUnauthorizedError("invalid token: rooming_house_id not found"))
			}

			// Convert the string IDs to uuid.UUID
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return utils.HandlerError(c, utils.NewUnauthorizedError("invalid token: invalid user_id"))
			}

			roomingHouseID, err := uuid.Parse(roomingHouseIDStr)
			if err != nil {
				return utils.HandlerError(c, utils.NewUnauthorizedError("invalid token: invalid rooming_house_id"))
			}

			// Role is still a string, so it's fine
			role := claims["role"].(string)

			c.Set("userPayload", &models.JWTPayload{
				UserID:         userID,
				Role:           role,
				RoomingHouseID: roomingHouseID,
			})
		} else {
			return utils.HandlerError(c, utils.NewUnauthorizedError("invalid token"))
		}

		return next(c)
	}
}
