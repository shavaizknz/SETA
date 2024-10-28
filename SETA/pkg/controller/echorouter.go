package controller

import (
	"net/http"
	"seta/pkg/model"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(handlers ...model.IController) *echo.Echo {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	})

	e.GET("/-/healthy", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "healthy",
		})
	})

	return e
}
