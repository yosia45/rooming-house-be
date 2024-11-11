package middlewares

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func CORSConfig() echo.MiddlewareFunc {
	origins := []string{"https://example.com"} // Default production
	if os.Getenv("APP_ENV") == "development" {
		origins = append(origins, "http://localhost:3000")
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	})
}
