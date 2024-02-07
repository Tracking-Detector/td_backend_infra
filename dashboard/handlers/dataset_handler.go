package handlers

import (
	"github.com/Tracking-Detector/td_backend_infra/dashboard/views/layouts"
	"github.com/gofiber/fiber/v2"
)

type DatasetHandler struct {
	app *fiber.App
}

func NewDatasetHandler(app *fiber.App) *DatasetHandler {
	return &DatasetHandler{
		app: app,
	}
}

func (h *DatasetHandler) Index(c *fiber.Ctx) error {
	return Render(c, layouts.Dashboard("Datasets"))
}

func (h *DatasetHandler) RegisterHandlers() {
	h.app.Get("/datasets", h.Index)
}
