package server

import "github.com/gofiber/fiber/v2"

func (h *Handler) registerRoutes(app *fiber.App) {

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/order/:order_uid", h.getOrderHandler)
}
