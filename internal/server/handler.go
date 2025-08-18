package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/merkulovlad/wbtech-go/internal/logger"
	"github.com/merkulovlad/wbtech-go/internal/model"
	ordr "github.com/merkulovlad/wbtech-go/internal/service/order"
)

type Handler struct {
	Order  ordr.Service
	Logger logger.InterfaceLogger
}

func NewHandler(order ordr.Service, logger logger.InterfaceLogger) *Handler {
	return &Handler{
		Order:  order,
		Logger: logger,
	}
}

func (h *Handler) getOrderHandler(c *fiber.Ctx) error {
	id := c.Params("order_uid")
	h.Logger.Infof("Getting order %s", id)
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&model.ErrorResponse{Status: fiber.StatusBadRequest, Msg: "Invalid id"})
	}
	order, err := h.Order.Get(c.Context(), id)
	if err != nil {
		h.Logger.Errorf("Get order error: %s", err.Error())
	}
	h.Logger.Infof("Get order %v", order)
	return c.Status(fiber.StatusOK).JSON(&order)
}
