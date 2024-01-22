package controller

import (
	"tds/shared/response"
	"tds/shared/service"

	"github.com/gofiber/fiber/v2"
)

type ExtractorController struct {
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
