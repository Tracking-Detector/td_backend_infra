package handlers

import (
	"fmt"

	"github.com/Tracking-Detector/td_backend_infra/dashboard/models"
	"github.com/Tracking-Detector/td_backend_infra/dashboard/services"
	"github.com/Tracking-Detector/td_backend_infra/dashboard/views/components"
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

func (h *DatasetHandler) Reload(c *fiber.Ctx) error {
	h.datasetService.LoadAllDatasets()
	return Render(c, dataset.Content(h.datasetService.GetAllDatasets()))
}

func (h *DatasetHandler) Create(c *fiber.Ctx) error {
	return Render(c, dataset.Create())
}

func (h *DatasetHandler) CreateDataset(c *fiber.Ctx) error {
	form := new(models.CreateDatasetPayload)
	if err := c.BodyParser(form); err != nil {
		return err
	}
	created, err := h.datasetService.CreateDataset(form)
	if err != nil {
		return Render(c, components.ErrorAlert("Error:", err.Error()))
	}
	return Render(c, components.InfoAlert("Success:", fmt.Sprintf("Dataset %s created", created.Name)))
}

func (h *DatasetHandler) Delete(c *fiber.Ctx) error {
	datset, err := h.datasetService.GetDatasetByID(c.Params("id"))
	if err != nil {
		return Render(c, components.ErrorAlert("Error:", err.Error()))
	}
	return Render(c, dataset.DeleteDialog(datset))
}

func (h *DatasetHandler) DeleteDataset(c *fiber.Ctx) error {
	err := h.datasetService.DeleteDataset(c.Params("id"))
	if err != nil {
		return Render(c, components.ErrorAlert("Error:", err.Error()))
	}
	return Render(c, components.InfoAlert("Success:", "Dataset deleted"))
}

func (h *DatasetHandler) RegisterHandlers() {
	h.app.Get("/datasets", h.Index)
	h.app.Get("/datasets/reload", h.Reload)
	h.app.Get("/datasets/create", h.Create)
	h.app.Post("/datasets/create", h.CreateDataset)
	h.app.Get("/datasets/delete/:id", h.Delete)
	h.app.Delete("/datasets/delete/:id", h.DeleteDataset)
}
