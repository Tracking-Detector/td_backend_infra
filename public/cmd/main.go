package main

import (
	"encoding/json"
	"fmt"

	"github.com/Tracking-Detector/td_backend_infra/public/handler"
	"github.com/Tracking-Detector/td_backend_infra/public/resources"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	res := resources.LoadHomeResource()
	t, _ := json.Marshal(res)
	fmt.Println(string(t))
	homeHandler := handler.NewHomeHandler(res)
	app.Static("/static", "static")
	app.GET("/", homeHandler.HandleHomeShow)

	app.Start(":8081")
}
