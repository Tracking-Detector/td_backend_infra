package controller

import (
	"net/http"
	"tds/shared/models"
	"tds/shared/response"
	"tds/shared/service"
	"tds/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type ModelController struct {
	app                *fiber.App
	trainingrunService service.ITrainingrunService
	modelService       service.IModelService
}

func NewModelController(trainingrunService service.ITrainingrunService, modelService service.IModelService) *ModelController {
	return &ModelController{
		trainingrunService: trainingrunService,
		modelService:       modelService,
	}
}

func (tc *ModelController) GetTrainingRuns(c *fiber.Ctx) error {
	runs, err := tc.trainingrunService.FindAllTrainingRuns(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(response.NewSuccessResponse(runs))
}

func (tc *ModelController) GetTrainingRunsByModelId(c *fiber.Ctx) error {
	modelId := c.Params("id")
	runs, err := tc.trainingrunService.FindAllByModelId(c.Context(), modelId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(response.NewSuccessResponse(runs))
}

func (tc *ModelController) GetAllModels(c *fiber.Ctx) error {
	models, err := tc.modelService.GetAllModels(c.Context())
	if err != nil {
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse(err.Error()))
		}
	}
	return c.Status(http.StatusOK).JSON(response.NewSuccessResponse(models))
}

func (tc *ModelController) CreateModel(c *fiber.Ctx) error {
	var model *models.Model
	if err := c.BodyParser(&model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.NewErrorResponse(err.Error()))
	}
	model, err := tc.modelService.Save(c.Context(), model)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(http.StatusCreated).JSON(response.NewSuccessResponse(model))

}

func (tc *ModelController) GetModelById(c *fiber.Ctx) error {
	modelId := c.Params("id")
	model, err := tc.modelService.GetModelById(c.Context(), modelId)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(err)
	}
	return c.Status(http.StatusOK).JSON(model)
}

func (tc *ModelController) DeleteModelById(c *fiber.Ctx) error {
	modelId := c.Params("id")
	if err := tc.modelService.DeleteModelByID(c.Context(), modelId); err != nil {
		return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(response.NewSuccessResponse("Successfully deleted model."))
}

func (tc *ModelController) DeleteRunById(c *fiber.Ctx) error {
	runId := c.Params("runId")
	if err := tc.trainingrunService.DeleteByID(c.Context(), runId); err != nil {
		return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(response.NewSuccessResponse("Successfully deleted model run."))
}

func (tc *ModelController) Start() {
	tc.app = fiber.New()
	tc.app.Use(cors.New())
	tc.app.Use(logger.New())
	tc.app.Get("/models/health", utils.GetHealth)
	tc.app.Get("/models", tc.GetAllModels)
	tc.app.Post("/models", tc.CreateModel)
	tc.app.Get("/models/:id", tc.GetModelById)
	tc.app.Delete("/models/:id", tc.DeleteModelById)
	tc.app.Get("/models/:id/runs", tc.GetTrainingRunsByModelId)
	tc.app.Delete("/models/:id/runs/:runId", tc.DeleteRunById)
	tc.app.Get("/models/runs", tc.GetTrainingRuns)
	tc.app.Listen(":8081")
}

func (tc *ModelController) Stop() {
	tc.app.Shutdown()
}
