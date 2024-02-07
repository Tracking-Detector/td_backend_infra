package handlers

import (
	"github.com/Tracking-Detector/td_backend_infra/dashboard/views/layouts"
	"github.com/gofiber/fiber/v2"
)

type ModelHandler struct {
	app *fiber.App
}

func NewModelHandler(app *fiber.App) *ModelHandler {
	return &ModelHandler{
		app: app,
	}
}

func (h *ModelHandler) Index(c *fiber.Ctx) error {
	return Render(c, layouts.Dashboard("Models"))
}

func (h *ModelHandler) RegisterHandlers() {
	h.app.Get("/models", h.Index)
}
