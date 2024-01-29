package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"tds/shared/configs"
	"tds/shared/repository"
	"tds/shared/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

type ProxyConfig struct {
	Path   string `json:"path"`
	Target string `json:"target"`
	Secure bool   `json:"secure"`
	Admin  bool   `json:"admin"`
}

func loadConfig(filePath string) ([]ProxyConfig, error) {
	var configs []ProxyConfig

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func main() {
	ctx := context.Background()
	db := configs.ConnectDB(ctx)
	userRepo := repository.NewMongoUserRepository(configs.GetDatabase(db))
	encryptionService := service.NewEncryptionService()
	authService := service.NewAuthService(userRepo, encryptionService)

	configs, err := loadConfig("proxy_config.json")
	if err != nil {
		log.Fatal("Error loading proxy configuration:", err)
	}

	app := fiber.New()

	for _, config := range configs {
		if config.Secure {
			app.Use(config.Path, func(c *fiber.Ctx) error {
				apiKey := c.Get("X-API-Key")
				res, err := authService.ValidateBearerToken(c.Context(), apiKey, config.Admin)
				if res && err == nil {
					return c.Next()
				} else {
					return c.SendStatus(http.StatusForbidden)
				}
			})
		}
		app.Use(config.Path, proxy.Balancer(proxy.Config{
			Servers: []string{
				config.Target,
			},
			ModifyRequest: func(c *fiber.Ctx) error {
				c.Request().SetRequestURI(config.Target + string(c.Request().RequestURI()))
				return nil
			},
			ModifyResponse: func(c *fiber.Ctx) error {
				c.Response().Header.Set("Access-Control-Allow-Origin", "*")
				return nil
			},
		}))
	}

	err = app.Listen(":8081")
	if err != nil {
		log.Fatal("Error starting Fiber server:", err)
	}
}
