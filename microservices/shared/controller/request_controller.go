package controller

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"tds/shared/models"
	"tds/shared/response"
	"tds/shared/service"
	"tds/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type RequestController struct {
	app            *fiber.App
	requestService service.IRequestService
}

func NewRequestController(requestService service.IRequestService) *RequestController {
	return &RequestController{
		requestService: requestService,
	}
}

func (rc *RequestController) CreateMultipleRequestData(c *fiber.Ctx) error {
	var requestData []*models.RequestData

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.NewErrorResponse(err.Error()))
	}
	err := rc.requestService.InsertManyRequests(c.Context(), requestData)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(http.StatusCreated).JSON(response.NewSuccessResponse("Successfully inserted requests."))
}

func (rc *RequestController) GetRequestById(c *fiber.Ctx) error {
	requestId := c.Params("id")
	request, err := rc.requestService.GetRequestById(c.Context(), requestId)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(response.NewSuccessResponse(request))
}

func (rc *RequestController) SearchRequests(c *fiber.Ctx) error {
	url := c.Query("url")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.NewErrorResponse("Query params 'page' and 'pageSize' should be left empty or numeric."))
	}
	requests, err := rc.requestService.GetPagedRequestsFilterdByUrl(c.Context(), url, page, pageSize)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.NewErrorResponse(err.Error()))
	}
	total, err := rc.requestService.CountDocumentsForUrlFilter(c.Context(), url)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse(err.Error()))
	}
	numPages := math.Ceil(float64(int(total)) / float64(pageSize))
	var next string
	if numPages < float64(page) {
		next = "/requests?page=" + fmt.Sprint(page+1) + "&pageSize=" + fmt.Sprint(pageSize) + "&url=" + url
	}
	return c.Status(http.StatusOK).JSON(response.NewPagedSuccessResponse(requests,
		"/requests?page="+fmt.Sprint(page)+"&pageSize="+fmt.Sprint(pageSize)+"&url="+url,
		next,
		int(numPages)))
}

func (rc *RequestController) CreateRequestData(c *fiber.Ctx) error {
	var requestData *models.RequestData
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.NewErrorResponse(err.Error()))
	}
	requestData, err := rc.requestService.SaveRequest(c.Context(), requestData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(http.StatusCreated).JSON(response.NewSuccessResponse(requestData))
}

func (rc *RequestController) Start() {
	rc.app = fiber.New()
	rc.app.Use(cors.New())
	rc.app.Use(logger.New())
	rc.app.Get("/requests/health", utils.GetHealth)
	rc.app.Get("/requests/:id", rc.GetRequestById)
	rc.app.Post("/requests", rc.CreateRequestData)
	rc.app.Post("/requests/multiple", rc.CreateMultipleRequestData)
	rc.app.Get("/requests", rc.SearchRequests)
	rc.app.Listen(":8081")
}

func (rc *RequestController) Stop() {
	rc.app.Shutdown()
}
