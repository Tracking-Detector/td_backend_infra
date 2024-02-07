package handlers

import (
	"github.com/Tracking-Detector/td_backend_infra/dashboard/views/layouts"
	"github.com/gofiber/fiber/v2"
)

type ExportHandler struct {
	app *fiber.App
}

func NewExportHandler(app *fiber.App) *ExportHandler {
	return &ExportHandler{
		app: app,
	}
}

func (h *ExportHandler) Index(c *fiber.Ctx) error {
	return Render(c, layouts.Dashboard("Exports"))
}

func (h *ExportHandler) RegisterHandlers() {
	h.app.Get("/exports", h.Index)
}
