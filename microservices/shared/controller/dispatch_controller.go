package controller

import (
	"net/http"
	"reflect"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/converter"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/response"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type DispatchController struct {
	app              *fiber.App
	exporterService  service.IExporterService
	exportRunService service.IExportRunService
	publishService   service.IPublishService
	modelService     service.IModelService
	datasetService   service.IDatasetService
}

func NewDispatchController(exporterService service.IExporterService, publishService service.IPublishService, modelService service.IModelService, datasetService service.IDatasetService, exportRunService service.IExportRunService) *DispatchController {
	return &DispatchController{
		exporterService:  exporterService,
		publishService:   publishService,
		modelService:     modelService,
		exportRunService: exportRunService,
		datasetService:   datasetService,
	}
}

func (dc *DispatchController) DispatchExportJob(c *fiber.Ctx) error {
	exporterId := c.Params("exporterId")
	reducer := c.Params("reducer")
	dataset := c.Params("dataset")
	if dataset == "all" {
		dataset = ""
	} else {
		if !dc.datasetService.IsValidDataset(c.Context(), dataset) {
			return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse("The dataset for the given id does not exist."))
		}
	}
	if !converter.IsValidReduceType(reducer) {
		return c.Status(http.StatusBadRequest).JSON(response.NewErrorResponse("The reducer type is not valid"))
	}
	isValid := dc.exporterService.IsValidExporter(c.Context(), exporterId)
	if !isValid {
		return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse("The extractor for the given id does not exist."))
	}
	dc.publishService.EnqueueExportJob(exporterId, reducer, dataset)
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
	if exist, err := dc.exportRunService.ExistByExporterIDAndRecducer(c.Context(), exporterId, reducer); err != nil || !exist {
		return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse("The export for the given id and reducer does not exist."))
	}
	model, err := dc.modelService.GetModelById(c.Context(), modelId)
	if err != nil || model == nil {
		return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse("The model for the given id does not exist."))
	}
	if !reflect.DeepEqual(exporter.Dimensions, model.Dims) {
		return c.Status(http.StatusBadRequest).JSON(response.NewErrorResponse("There is a dimension mismatch for the dataset and the model."))
	}

	dc.publishService.EnqueueTrainingJob(modelId, exporterId, reducer)

	return c.Status(http.StatusCreated).JSON(response.NewSuccessResponse("The training job has been dispatched."))
}

func (dc *DispatchController) Start() {
	dc.app = fiber.New()
	dc.app.Use(cors.New())
	dc.app.Use(logger.New())
	dc.app.Get("/dispatch/health", utils.GetHealth)
	dc.app.Post("/dispatch/export/:exporterId/:reducer/:dataset", dc.DispatchExportJob)
	dc.app.Post("/dispatch/train/:modelId/run/:exporterId/:reducer", dc.DispatchTrainingJob)

	dc.app.Listen(":8081")
}

func (dc *DispatchController) Stop() {
	dc.app.Shutdown()
}
