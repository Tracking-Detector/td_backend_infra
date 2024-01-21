package responses

import (
	"tds/shared/models"

	"github.com/gofiber/fiber/v2"
)

type RequestDataResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}

type PagedRequestDataResponse struct {
	Self     string               `json:"_self"`
	Next     string               `json:"_next"`
	PageSize int                  `json:"pageSize"`
	NumPages float64              `json:"numPages"`
	Count    int                  `json:"count"`
	Content  []models.RequestData `json:"content"`
}
