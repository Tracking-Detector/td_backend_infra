package controller

import (
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/payload"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/response"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type DatasetController struct {
	app            *fiber.App
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

func (dc *DatasetController) Start() {
	dc.app = fiber.New()
	dc.app.Use(cors.New())
	dc.app.Use(logger.New())
	dc.app.Get("/datasets/health", utils.GetHealth)
	dc.app.Get("/datasets", dc.GetAllDatasets)
	dc.app.Post("/datasets", dc.CreateDataset)
	dc.app.Listen(":8081")
}

func (dc *DatasetController) Stop() {
	dc.app.Shutdown()
}
