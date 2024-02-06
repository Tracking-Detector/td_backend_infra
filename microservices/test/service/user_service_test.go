package service_test

import (
	"context"
	"errors"
	"os"

	"testing"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
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
	user := &models.UserData{
		ID: "1", Email: "user1@example.com", Role: models.CLIENT,
	}
	var transactionError error
	suite.userRepo.On("InTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
		f := args.Get(1).(func(ctx context.Context) error)
		transactionError = f(context.Background())
	}).Return(nil)
	suite.userRepo.On("FindByID", mock.Anything, "1").Return(user, nil)
	suite.userRepo.On("DeleteByID", mock.Anything, "1").Return(nil)
	// when
	_ = suite.userService.DeleteUserByID(context.Background(), "1")

	// then
	suite.userRepo.AssertCalled(suite.T(), "InTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error"))
	suite.userRepo.AssertCalled(suite.T(), "FindByID", mock.Anything, "1")
	suite.userRepo.AssertCalled(suite.T(), "DeleteByID", mock.Anything, "1")
	suite.NoError(transactionError)
}

func (suite *UserServiceTest) TestDeleteUserByID_Failure_AdminUser() {
	// given
	user := &models.UserData{
		ID: "1", Email: "user1@example.com", Role: models.ADMIN,
	}
	var transactionError error
	suite.userRepo.On("FindByID", mock.Anything, "1").Return(user, nil)
	suite.userRepo.On("InTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
		f := args.Get(1).(func(ctx context.Context) error)
		transactionError = f(context.Background())
	}).Return(nil)

	// when
	_ = suite.userService.DeleteUserByID(context.Background(), "1")

	// then
	suite.userRepo.AssertNotCalled(suite.T(), "DeleteByID")
	suite.Error(transactionError)
	suite.Contains(transactionError.Error(), "cannot delete admin users")
}

func (suite *UserServiceTest) TestDeleteUserByID_Failure_NoUserForID() {
	// given
	var transactionError error
	suite.userRepo.On("FindByID", mock.Anything, "1").Return(nil, errors.New("No user found."))
	suite.userRepo.On("InTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
		f := args.Get(1).(func(ctx context.Context) error)
		transactionError = f(context.Background())
	}).Return(nil)

	// when
	_ = suite.userService.DeleteUserByID(context.Background(), "1")

	// then
	suite.userRepo.AssertNotCalled(suite.T(), "DeleteByID")
	suite.Error(transactionError)
	suite.Contains(transactionError.Error(), "No user found.")
}

func (suite *UserServiceTest) TestInitAdmin_Success() {
	// given
	apiKey := "asdasdfasdfqwdasdasfadsfasdf"
	adminEmail := "test@test.com"
	os.Setenv("ADMIN_API_KEY", apiKey)
	os.Setenv("EMAIL", adminEmail)
	suite.encryptionService.On("HashPassword", apiKey).Return("digest", nil)
	user := &models.UserData{
		Role:  models.ADMIN,
		Email: adminEmail,
		Key:   "digest",
	}
	suite.userRepo.On("Save", mock.Anything, user).Return(user, nil)

	// when
	err := suite.userService.InitAdmin(context.Background())
	// then
	suite.NoError(err)
	suite.encryptionService.AssertCalled(suite.T(), "HashPassword", apiKey)
	suite.userRepo.AssertCalled(suite.T(), "Save", mock.Anything, user)
}

func (suite *UserServiceTest) TestCreateApiUser_Success() {
	// given
	email := "test@test.com"
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
	pw, _ := suite.userService.CreateApiUser(context.Background(), email)
	// then
	suite.Equal(pw, "pw")
	suite.NoError(transactionError)
}

func (suite *UserServiceTest) TestCreateApiUser_ErrorEmailAlreadyExists() {
	// given
	email := "test@test.com"
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
	pw, _ := suite.userService.CreateApiUser(context.Background(), email)
	// then
	suite.Empty(pw)
	suite.Contains(transactionError.Error(), "user with email already registered")
}
