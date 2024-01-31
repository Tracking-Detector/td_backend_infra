package controller

import (
	"tds/shared/payload"
	"tds/shared/response"
	"tds/shared/service"

	"github.com/gofiber/fiber/v2"
)

type DatasetController struct {
	datasetService service.IDatasetService
}

func NewDatasetController(datasetService service.IDatasetService) *DatasetController {
	return &DatasetController{
		datasetService: datasetService,
	}
}

func (dc *DatasetController) GetAllDatasets(c *fiber.Ctx) error {
	datasets := dc.datasetService.GetAllDatasets()
	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(datasets))
}

func (dc *DatasetController) CreateDataset(c *fiber.Ctx) error {
	datasetPayload := new(payload.CreateDatasetPayload)
	if err := c.BodyParser(datasetPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(err.Error()))
	}
	dataset, err := dc.datasetService.CreateDataset(c.Context(), datasetPayload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(fiber.StatusCreated).JSON(response.NewSuccessResponse(dataset))
}
