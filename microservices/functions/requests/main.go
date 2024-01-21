package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"tds/shared/configs"
	"tds/shared/models"
	"tds/shared/responses"
	"tds/shared/utils"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var logger = configs.GetLogger("request-api")
var requestDataCollection *mongo.Collection = configs.GetCollection(configs.DB, configs.EnvRequestCollection())
var validate = validator.New()

func SearchRequests(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var requestDataValues []models.RequestData
	findOptions := options.Find()
	filter := bson.M{}
	url := c.Query("url")
	if url != "" {
		filter = bson.M{
			"url": bson.M{
				"$regex": primitive.Regex{
					Pattern: url,
					Options: "i",
				},
			},
		}
	}
	total, _ := requestDataCollection.CountDocuments(ctx, filter)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	findOptions.SetSkip((int64(page) - 1) * int64(pageSize))
	findOptions.SetLimit(int64(pageSize))
	cursor, _ := requestDataCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var requestData models.RequestData
		cursor.Decode(&requestData)
		requestDataValues = append(requestDataValues, requestData)
	}
	numPages := math.Ceil(float64(int(total)) / float64(pageSize))
	var next string
	if numPages < float64(page) {
		next = "/requests?page=" + fmt.Sprint(page+1) + "&pageSize=" + fmt.Sprint(pageSize) + "&url=" + url
	}
	return c.Status(http.StatusOK).JSON(responses.PagedRequestDataResponse{
		Self:     "/requests?page=" + fmt.Sprint(page) + "&pageSize=" + fmt.Sprint(pageSize) + "&url=" + url,
		Next:     next,
		PageSize: pageSize,
		NumPages: numPages,
		Count:    int(total),
		Content:  requestDataValues,
	})
}

func GetRequestDataById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	requestId := c.Params("requestId")
	var requestData models.RequestData
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(requestId)

	err := requestDataCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&requestData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestDataResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.RequestDataResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": requestData}})
}

func CreateRequestData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var requestData models.RequestData
	defer cancel()

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestDataResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&requestData); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestDataResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": validationErr.Error()}})
	}

	newUser := models.RequestData{
		Id:                primitive.NewObjectID(),
		DocumentId:        requestData.DocumentId,
		DocumentLifecycle: requestData.DocumentLifecycle,
		FrameId:           requestData.FrameId,
		FrameType:         requestData.FrameType,
		Initiator:         requestData.Initiator,
		Method:            requestData.Method,
		ParentFrameId:     requestData.ParentFrameId,
		RequestId:         requestData.RequestId,
		TabId:             requestData.TabId,
		TimeStamp:         requestData.TimeStamp,
		Type:              requestData.Type,
		URL:               requestData.URL,
		RequestHeaders:    requestData.RequestHeaders,
		Response:          requestData.Response,
		Success:           requestData.Success,
		Labels:            requestData.Labels,
	}

	result, err := requestDataCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestDataResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.RequestDataResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"data": result}})
}

func CreateMultipleRequestData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var requestData []models.RequestData
	defer cancel()

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestDataResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()}})
	}
	var interfaceSlice []interface{} = make([]interface{}, len(requestData))
	for i, d := range requestData {
		interfaceSlice[i] = d
	}
	_, err := requestDataCollection.InsertMany(ctx, interfaceSlice)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestDataResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.RequestDataResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"data": "Successfully inserted requests."}})
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/requests/health", utils.GetHealth)
	app.Get("/requests/:requestId", GetRequestDataById)
	app.Post("/requests", CreateRequestData)
	app.Post("/requests/multiple", CreateMultipleRequestData)
	app.Get("/requests", SearchRequests)
	// logger.Info("Server Running")
	app.Listen(":8081")
}
