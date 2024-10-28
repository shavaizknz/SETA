package main

import (
	"context"
	"log"
	"seta/pkg/clients/paymentgateway"
	"seta/pkg/clients/paymentgateway/paymentgatewaya"
	"seta/pkg/clients/paymentgateway/paymentgatewayb"
	"seta/pkg/config"
	"seta/pkg/controller"
	"seta/pkg/infra/pg"
	"seta/pkg/logger"
	"seta/pkg/repository"
	"seta/pkg/service"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title SETA API
// @version 1.0
// @description This is a payment gateway API doc for SETA

// @host localhost:8080/
// @Schemes http
func main() {
	godotenv.Load()
	config := config.GetConfigManager()

	dbPool, err := pg.DBPoolProvider(config.GetDatabaseDSN(), context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer dbPool.DB.Close()

	transactionService := service.TransactionServiceProvider(repository.TransactionRepositoryProvider(dbPool.DB), []paymentgateway.IPaymentGateway{
		paymentgatewaya.ClientProvider(config.GetGatewayAEndpoint()),
		paymentgatewayb.ClientProvider(config.GetGatewayBEndpoint()),
	})

	transactionController := controller.TransactionControllerProvider(transactionService)

	e := controller.SetupRoutes(transactionController)

	transactionController.SetupRoutes(e.Group("/api/v1"))
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(logger.LogMiddleware)
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "swagger") {
				return true
			}
			return false
		},
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.Start(":8080")
}
