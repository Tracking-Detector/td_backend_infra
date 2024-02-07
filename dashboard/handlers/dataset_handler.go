package handlers

import (
	"github.com/Tracking-Detector/td_backend_infra/dashboard/services"
	"github.com/Tracking-Detector/td_backend_infra/dashboard/views/pages/dataset"
	"github.com/gofiber/fiber/v2"
)

type DatasetHandler struct {
	app            *fiber.App
	datasetService services.IDatasetService
}

func NewDatasetHandler(app *fiber.App, datasetService services.IDatasetService) *DatasetHandler {
	return &DatasetHandler{
		app:            app,
		datasetService: datasetService,
	}
}

func (h *DatasetHandler) Index(c *fiber.Ctx) error {
	return Render(c, dataset.Index(h.datasetService.GetAllDatasets()))
}

func (h *DatasetHandler) Create(c *fiber.Ctx) error {
	return Render(c, dataset.Create())
}

func (h *DatasetHandler) RegisterHandlers() {
	h.app.Get("/datasets", h.Index)
	h.app.Get("/datasets/create", h.Create)
}
