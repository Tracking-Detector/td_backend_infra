package controller

import (
	"net/http"
	"tds/shared/response"
	"tds/shared/service"

	"github.com/gofiber/fiber/v2"
)

type TrainingrunController struct {
	trainingrunService service.ITrainingrunService
}

func NewTrainingrunController(trainingrunService service.ITrainingrunService) *TrainingrunController {
	return &TrainingrunController{
		trainingrunService: trainingrunService,
	}
}

func (tc *TrainingrunController) GetTrainingRuns(c *fiber.Ctx) error {
	runs, err := tc.trainingrunService.FindAllTrainingRuns(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(response.NewSuccessResponse(runs))
}

func (tc *TrainingrunController) GetTrainingRunsByModelName(c *fiber.Ctx) error {
	modelName := c.Params("modelName")
	runs, err := tc.trainingrunService.FindAllTrainingRunsForModelname(c.Context(), modelName)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(response.NewSuccessResponse(runs))
}
