package handler

import (
	"github.com/Tracking-Detector/td_backend_infra/public/views/home"
	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
}

func (h *HomeHandler) HandleHomeShow(c echo.Context) error {
	return render(c, home.Show())
}
