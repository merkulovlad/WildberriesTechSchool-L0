package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/merkulovlad/wbtech-go/internal/logger"
	"github.com/merkulovlad/wbtech-go/internal/service/order"
)

func NewServer(orderSvc order.Service, log logger.InterfaceLogger) *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001",
		AllowMethods:     "GET,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: false,
	}))
	h := NewHandler(orderSvc, log)
	h.registerRoutes(app)

	return app
}
