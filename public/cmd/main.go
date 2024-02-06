package main

import (
	"github.com/Tracking-Detector/td_backend_infra/public/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	homeHandler := &handler.HomeHandler{}

	app.GET("/", homeHandler.HandleHomeShow)

	app.Start(":8081")
}
