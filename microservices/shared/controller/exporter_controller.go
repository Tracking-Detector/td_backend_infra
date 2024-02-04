package controller

import (
	"tds/shared/response"
	"tds/shared/service"
	"tds/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type ExtractorController struct {
	app              *fiber.App
	extractorService service.IExporterService
}

func NewExtractorController(extractorService service.IExporterService) *ExtractorController {
	return &ExtractorController{
		extractorService: extractorService,
	}
}

func (con *ExtractorController) GetAllExporter(c *fiber.Ctx) error {
	extractors, err := con.extractorService.GetAllExporter(c.Context())
	if err != nil {
		errorResponse := response.NewErrorResponse(err.Error())
		return c.Status(500).JSON(errorResponse)
	}
	successResponse := response.NewSuccessResponse(extractors)
	return c.Status(200).JSON(successResponse)
}

func (con *ExtractorController) Start() {
	con.app = fiber.New()
	con.app.Use(cors.New())
	con.app.Use(logger.New())
	con.app.Get("/export/health", utils.GetHealth)
	con.app.Get("/export", con.GetAllExporter)
	con.app.Listen(":8081")
}

func (con *ExtractorController) Stop() {
	con.app.Shutdown()
}
