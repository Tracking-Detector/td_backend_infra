package handlers

import (
	"github.com/Tracking-Detector/td_backend_infra/dashboard/services"
	"github.com/Tracking-Detector/td_backend_infra/dashboard/views/pages/home"
	"github.com/gofiber/fiber/v2"
)

type IHandler interface {
	RegisterHandlers()
}

type HomeHandler struct {
	app           *fiber.App
	statusService services.IStatusService
}

func NewHomeHandler(app *fiber.App, statusService services.IStatusService) *HomeHandler {
	return &HomeHandler{
		app:           app,
		statusService: statusService,
	}
}

func (h *HomeHandler) Index(c *fiber.Ctx) error {
	return Render(c, home.Index(h.statusService.GetStatus()))
}

func (h *HomeHandler) ReloadContent(c *fiber.Ctx) error {
	h.statusService.ReloadCache()
	return Render(c, home.IndexContent(h.statusService.GetStatus()))
}

func (h *HomeHandler) RegisterHandlers() {
	h.app.Get("/", h.Index)
	h.app.Get("/reload", h.ReloadContent)
}
