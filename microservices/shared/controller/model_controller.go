package controller

import (
	"net/http"
	"tds/shared/models"
	"tds/shared/response"
	"tds/shared/service"

	"github.com/gofiber/fiber/v2"
)

type ModelController struct {
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
	if err := tc.modelService.Save(c.Context(), model); err != nil {
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
