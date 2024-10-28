package model

import "github.com/labstack/echo/v4"

type DefaultResponse struct {
	Data interface{} `json:"data"`
}

type DefaultError struct {
	Error string `json:"error"`
}

type IController interface {
	SetupRoutes(r *echo.Group)
}
