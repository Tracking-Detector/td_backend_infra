package main

import (
	"context"

	"net/http"
	"tds/shared/configs"
	"tds/shared/models"
	"tds/shared/responses"
	"tds/shared/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func InitAdmin() {
	adminApiKey := configs.EnvAdminApiKey()
	hash, err := bcrypt.GenerateFromPassword([]byte(adminApiKey), bcrypt.DefaultCost)

	if err != nil {
		log.WithFields(log.Fields{
			"service": "users",
			"error":   err.Error(),
		}).Fatal("Could not generate Hash for Admin key.")
	}

	userCollection := configs.GetCollection(configs.ConnectDB(), "users")

	result := userCollection.FindOne(context.Background(), bson.M{
		"role": "admin",
	})
	var adminUser models.UserData
	if err := result.Decode(&adminUser); err != nil {
		admin := models.UserData{
			Role:  models.ADMIN,
			Email: "henry.schwerdtner@web.de",
			Key:   string(hash),
		}
		userCollection.InsertOne(context.Background(), admin)
		return
	}

	comparingResult := bcrypt.CompareHashAndPassword([]byte(adminUser.Key), []byte(adminApiKey))
	if comparingResult != nil {
		userCollection.UpdateOne(context.Background(), bson.M{
			"role": "admin",
		}, bson.M{
			"$set": bson.M{
				"key": string(hash),
			},
		})
	}
}

func CreateAPIKey(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var userData models.CreateUserData
	defer cancel()

	if err := c.BodyParser(&userData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestDataResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()}})
	}
	userCollection := configs.GetCollection(configs.ConnectDB(), "users")
	result := userCollection.FindOne(ctx, bson.M{
		"email": userData.Email,
	})
	if result.Err() != mongo.ErrNoDocuments {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestDataResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": "user already exisits."}})
	}
	key := uuid.New().String()
	hash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestDataResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": "Could not generate hash of api key."}})
	}
	newUser := models.UserData{
		Role:  models.CLIENT,
		Email: userData.Email,
		Key:   string(hash),
	}
	userCollection.InsertOne(ctx, newUser)
	return c.Status(http.StatusCreated).JSON(responses.RequestDataResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"data": "User created with API-Key: '" + key + "'"}})

}

func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection := configs.GetCollection(configs.ConnectDB(), "users")
	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestDataResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err}})
	}
	var userData []models.UserDataRepresentation
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user models.UserDataRepresentation
		cursor.Decode(&user)
		userData = append(userData, user)
	}
	return c.Status(http.StatusOK).JSON(responses.RequestDataResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": userData}})
}

func DeleteUserById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userId := c.Params("userId")
	userCollection := configs.GetCollection(configs.ConnectDB(), "users")
	objectId, objectIdError := primitive.ObjectIDFromHex(userId)
	if objectIdError != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestDataResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": "No objectId passed as a param."}})
	}
	result := userCollection.FindOne(ctx, bson.M{"_id": objectId})
	var userToDelete models.UserData

	if err := result.Decode(&userToDelete); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestDataResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()}})
	}

	if userToDelete.Role == models.ADMIN {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestDataResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": "Impossible to delete the system admin user."}})
	}

	_, deletionError := userCollection.DeleteOne(ctx, bson.M{
		"_id":   objectId,
		"email": userToDelete.Email,
	})

	if deletionError != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestDataResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": deletionError.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.RequestDataResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": "User has been deleted."}})

}

func main() {
	InitAdmin()
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/users/health", utils.GetHealth)
	app.Get("/users", GetUsers)
	app.Post("/users", CreateAPIKey)
	app.Delete("/users/:userId", DeleteUserById)
	app.Listen(":8081")
}
