package controller

import (
	"net/http"
	"reflect"
	"tds/shared/converter"
	"tds/shared/response"
	"tds/shared/service"
	"tds/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type DispatchController struct {
	app             *fiber.App
	exporterService service.IExporterService
	publishService  service.IPublishService
	modelService    service.IModelService
}

func NewDispatchController(exporterService service.IExporterService, publishService service.IPublishService, modelService service.IModelService) *DispatchController {
	return &DispatchController{
		exporterService: exporterService,
		publishService:  publishService,
		modelService:    modelService,
	}
}

func (dc *DispatchController) DispatchExportJob(c *fiber.Ctx) error {
	exporterId := c.Params("exporterId")
	reducer := c.Params("reducer")
	if !converter.IsValidReduceType(reducer) {
		return c.Status(http.StatusBadRequest).JSON(response.NewErrorResponse("The reducer type is not valid"))
	}
	isValid := dc.exporterService.IsValidExporter(c.Context(), exporterId)
	if !isValid {
		return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse("The extractor for the given id does not exist."))
	}
	dc.publishService.EnqueueExportJob(exporterId, reducer, "") // TODO make dataset an option
	return c.Status(http.StatusCreated).JSON(response.NewSuccessResponse("The export has been dispatched."))
}

func (dc *DispatchController) DispatchTrainingJob(c *fiber.Ctx) error {
	modelId := c.Params("modelId")
	exporterId := c.Params("exporterId")
	reducer := c.Params("reducer")

	if !converter.IsValidReduceType(reducer) {
		return c.Status(http.StatusBadRequest).JSON(response.NewErrorResponse("The reducer type is not valid"))
	}
	exporter, err := dc.exporterService.FindByID(c.Context(), exporterId)
	if err != nil || exporter == nil {
		return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse("The extractor for the given id does not exist."))
	}
	model, err := dc.modelService.GetModelById(c.Context(), modelId)
	if err != nil || model == nil {
		return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse("The model for the given id does not exist."))
	}
	if !reflect.DeepEqual(exporter.Dimensions, model.Dims) {
		c.Status(http.StatusBadRequest).JSON(response.NewErrorResponse("There is a dimension mismatch for the dataset and the model."))
	}

	dc.publishService.EnqueueTrainingJob(modelId, exporterId, reducer)

	return c.Status(http.StatusCreated).JSON(response.NewSuccessResponse("The training job has been dispatched."))
}

func (dc *DispatchController) Start() {
	dc.app = fiber.New()
	dc.app.Use(cors.New())
	dc.app.Use(logger.New())
	dc.app.Get("/dispatch/health", utils.GetHealth)
	dc.app.Post("/dispatch/export/:exporterId/:reducer", dc.DispatchExportJob)
	dc.app.Post("/dispatch/train/:modelId/run/:exporterId/:reducer", dc.DispatchTrainingJob)

	dc.app.Listen(":8081")
}

func (dc *DispatchController) Stop() {
	dc.app.Shutdown()
}
