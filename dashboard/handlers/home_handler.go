package handlers

import (
	"github.com/Tracking-Detector/td_backend_infra/dashboard/views/layouts"
	"github.com/gofiber/fiber/v2"
)

type IHandler interface {
	RegisterHandlers()
}

type HomeHandler struct {
	app *fiber.App
}

func NewHomeHandler(app *fiber.App) *HomeHandler {
	return &HomeHandler{
		app: app,
	}
}

func (h *HomeHandler) Index(c *fiber.Ctx) error {
	return Render(c, layouts.Page("Home"))
}

func (h *HomeHandler) RegisterHandlers() {
	h.app.Get("/", h.Index)
}
