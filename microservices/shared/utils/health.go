package utils

import (
	"net/http"
	"tds/shared/response"

	"github.com/gofiber/fiber/v2"
)

func GetHealth(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(response.NewSuccessResponse("System is running correct."))
}
