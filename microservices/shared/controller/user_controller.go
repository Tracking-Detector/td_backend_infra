package controller

import (
	"net/http"
	"tds/shared/payload"
	"tds/shared/representation"
	"tds/shared/response"
	"tds/shared/service"
	"tds/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type UserController struct {
	app         *fiber.App
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) GetUsers(c *fiber.Ctx) error {
	users, err := uc.userService.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(representation.ConvertUserDatasToUserDataRepresentations(users))
}

func (uc *UserController) CreateApiUser(c *fiber.Ctx) error {
	var createUserData payload.CreateUserData
	if err := c.BodyParser(&createUserData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.NewErrorResponse(err.Error()))
	}
	key, err := uc.userService.CreateApiUser(c.Context(), createUserData.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse(err.Error()))
	}
	return c.Status(http.StatusCreated).JSON("User created with key '" + key + "'.")
}

func (uc *UserController) DeleteUserByID(c *fiber.Ctx) error {
	userId := c.Params("id")

	err := uc.userService.DeleteUserByID(c.Context(), userId)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(response.NewSuccessResponse("Deleted user successful."))
}

func (uc *UserController) Start() {
	uc.app = fiber.New()
	uc.app.Use(cors.New())
	uc.app.Use(logger.New())
	uc.app.Get("/users/health", utils.GetHealth)
	uc.app.Get("/users", uc.GetUsers)
	uc.app.Post("/users", uc.CreateApiUser)
	uc.app.Delete("/users/:Id", uc.DeleteUserByID)
	uc.app.Listen(":8081")
}

func (uc *UserController) Stop() {
	uc.app.Shutdown()
}
