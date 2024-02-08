package handler

import (
	"github.com/Tracking-Detector/td_backend_infra/public/models"
	"github.com/Tracking-Detector/td_backend_infra/public/views/home"
	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
	Home *models.Home
}

func NewHomeHandler(home *models.Home) *HomeHandler {
	return &HomeHandler{
		Home: home,
	}
}

func (h *HomeHandler) HandleHomeShow(c echo.Context) error {
	return render(c, home.Index(h.Home))
}
