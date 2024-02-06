package main

import (
	"github.com/Tracking-Detector/td_backend_infra/public/handler"
	"github.com/Tracking-Detector/td_backend_infra/public/resources"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	res := resources.LoadHomeResource()
	homeHandler := handler.NewHomeHandler(res)
	app.Static("/static", "static")
	app.GET("/", homeHandler.HandleHomeShow)

	app.Start(":8081")
}
