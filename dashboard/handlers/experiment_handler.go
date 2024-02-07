package handlers

import (
	"github.com/Tracking-Detector/td_backend_infra/dashboard/views/layouts"
	"github.com/gofiber/fiber/v2"
)

type ExperimentHandler struct {
	app *fiber.App
}

func NewExperimentHandler(app *fiber.App) *ExperimentHandler {
	return &ExperimentHandler{
		app: app,
	}
}

func (h *ExperimentHandler) Index(c *fiber.Ctx) error {
	return Render(c, layouts.Dashboard("Experiments"))
}

func (h *ExperimentHandler) RegisterHandlers() {
	h.app.Get("/experiments", h.Index)
}
