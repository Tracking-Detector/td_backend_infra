package responses

import "github.com/gofiber/fiber/v2"

type ExportTypesResponse struct {
	Status int          `json:"status"`
	Data   []*fiber.Map `json:"data"`
}

type ExportJobStartResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
