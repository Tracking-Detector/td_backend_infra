package service_test

import (
	"context"
	"errors"

	"testing"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/payload"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
	"github.com/Tracking-Detector/td_backend_infra/microservices/test/testsupport/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestUserServiceService(t *testing.T) {
	suite.Run(t, &UserServiceTest{})
}

type UserServiceTest struct {
	suite.Suite
	userService       *service.UserService
	userRepo          *mocks.UserRepository
	encryptionService *mocks.IEncryptionService
}

func (suite *UserServiceTest) SetupTest() {
	suite.encryptionService = new(mocks.IEncryptionService)
	suite.userRepo = new(mocks.UserRepository)
	suite.userService = service.NewUserService(suite.userRepo, suite.encryptionService)
}

func (suite *UserServiceTest) TestGetAllUsers_Success() {
	// given
	users := []*models.UserData{
		{ID: "1", Email: "user1@example.com", Role: models.CLIENT},
		{ID: "2", Email: "user2@example.com", Role: models.CLIENT},
	}
	suite.userRepo.On("FindAll", mock.Anything).Return(users, nil)

	// when
	result, err := suite.userService.GetAllUsers(context.Background())

	// then
	suite.userRepo.AssertCalled(suite.T(), "FindAll", mock.Anything)
	suite.NoError(err)
	suite.Len(result, 2)
}

func (suite *UserServiceTest) TestDeleteUserByID_Success() {
	// given

	suite.userRepo.On("DeleteByID", mock.Anything, "1").Return(nil)
	// when
	err := suite.userService.DeleteUserByID(context.Background(), "1")

	// then

	suite.userRepo.AssertCalled(suite.T(), "DeleteByID", mock.Anything, "1")
	suite.NoError(err)
}

func (suite *UserServiceTest) TestCreateApiUser_Success() {
	// given
	email := "test@test.com"
	userPayload := payload.CreateUserData{
		Role:  models.CLIENT,
		Email: email,
	}
	user := &models.UserData{
		Role:  models.CLIENT,
		Email: email,
		Key:   "digest",
	}
	suite.encryptionService.On("GenerateApiKey").Return("digest", "pw", nil)
	var transactionError error
	suite.userRepo.On("InTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
		f := args.Get(1).(func(ctx context.Context) error)
		transactionError = f(context.Background())
	}).Return(nil)
	suite.userRepo.On("FindByEmail", mock.Anything, email).Return(nil, errors.New("Error"))
	suite.userRepo.On("Save", mock.Anything, user).Return(user, nil)
	// when
	pw, _ := suite.userService.CreateApiUser(context.Background(), userPayload)
	// then
	suite.Equal(pw, "pw")
	suite.NoError(transactionError)
}

func (suite *UserServiceTest) TestCreateApiUser_ErrorEmailAlreadyExists() {
	// given
	email := "test@test.com"
	userPayload := payload.CreateUserData{
		Role:  models.CLIENT,
		Email: email,
	}
	user := &models.UserData{
		Role:  models.CLIENT,
		Email: email,
		Key:   "digest",
	}
	suite.encryptionService.On("GenerateApiKey").Return("digest", "pw", nil)
	var transactionError error
	suite.userRepo.On("InTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
		f := args.Get(1).(func(ctx context.Context) error)
		transactionError = f(context.Background())
	}).Return(errors.New("err"))
	suite.userRepo.On("FindByEmail", mock.Anything, email).Return(user, nil)
	// when
	pw, _ := suite.userService.CreateApiUser(context.Background(), userPayload)
	// then
	suite.Empty(pw)
	suite.Contains(transactionError.Error(), "user with email already registered")
}
