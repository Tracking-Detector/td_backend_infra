package utils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetHealth(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(bson.M{
		"status":  200,
		"message": "System is running correct.",
	})
}
