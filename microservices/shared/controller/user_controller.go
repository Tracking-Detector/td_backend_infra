package controller

import (
	"net/http"
	"tds/shared/payload"
	"tds/shared/representation"
	"tds/shared/response"
	"tds/shared/service"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
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
