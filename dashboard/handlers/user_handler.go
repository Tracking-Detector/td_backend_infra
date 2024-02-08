package handlers

import (
	"github.com/Tracking-Detector/td_backend_infra/dashboard/views/layouts"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	app *fiber.App
}

func NewUserHandler(app *fiber.App) *UserHandler {
	return &UserHandler{
		app: app,
	}
}

func (h *UserHandler) Index(c *fiber.Ctx) error {
	return Render(c, layouts.Dashboard("Users"))
}

func (h *UserHandler) RegisterHandlers() {
	h.app.Get("/users", h.Index)
}
