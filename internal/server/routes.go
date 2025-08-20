package server

import "github.com/gofiber/fiber/v2"

// Health check endpoint
// @Summary      Health check
// @Description  Returns service health status
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /healthz [get]
func (h *Handler) registerRoutes(app *fiber.App) {

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/order/:order_uid", h.getOrderHandler)
}
