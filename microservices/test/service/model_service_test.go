package service_test

import (
	"context"
	"errors"
	"tds/shared/models"
	"tds/shared/service"
	"tds/test/testsupport/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestModelService(t *testing.T) {
	suite.Run(t, &ModelServiceTest{})
}

type ModelServiceTest struct {
	suite.Suite
	modelService       *service.ModelService
	trainingRunService *mocks.ITrainingrunService
	modelRepo          *mocks.ModelRepository
}

func (suite *ModelServiceTest) SetupTest() {
	suite.trainingRunService = new(mocks.ITrainingrunService)
	suite.modelRepo = new(mocks.ModelRepository)
	suite.modelService = service.NewModelService(suite.modelRepo, suite.trainingRunService)
}

func (suite *ModelServiceTest) TestSave_Success() {
	// given
	model := &models.Model{
		ID:          "1",
		Name:        "ModelName",
		Description: "ModelDescription",
		Dims:        []int{204, 1},
	}
	suite.modelRepo.On("Save", mock.Anything, model).Return(nil)
	// when
	_, err := suite.modelService.Save(context.Background(), model)
	// then

	suite.NoError(err)
	suite.modelRepo.AssertCalled(suite.T(), "Save", mock.Anything, model)
}

func (suite *ModelServiceTest) TestSave_Error() {
	// given
	model := &models.Model{
		ID:          "1",
		Name:        "ModelName",
		Description: "ModelDescription",
		Dims:        []int{204, 1},
	}
	suite.modelRepo.On("Save", mock.Anything, model).Return(errors.New("error"))
	// when
	_, err := suite.modelService.Save(context.Background(), model)
	// then
	suite.Error(err)
	suite.modelRepo.AssertCalled(suite.T(), "Save", mock.Anything, model)
}

func (suite *ModelServiceTest) TestGetAllModels_Success() {
	// given
	models := []*models.Model{{
		ID:          "1",
		Name:        "ModelName",
		Description: "ModelDescription",
		Dims:        []int{204, 1},
	}}
	suite.modelRepo.On("FindAll", mock.Anything).Return(models, nil)
	// when
	_, err := suite.modelService.GetAllModels(context.Background())
	// then
	suite.NoError(err)
	suite.modelRepo.AssertCalled(suite.T(), "FindAll", mock.Anything)
}

func (suite *ModelServiceTest) TestGetAllModels_Error() {
	// given

	suite.modelRepo.On("FindAll", mock.Anything).Return(nil, errors.New("error"))
	// when
	_, err := suite.modelService.GetAllModels(context.Background())
	// then
	suite.Error(err)
	suite.modelRepo.AssertCalled(suite.T(), "FindAll", mock.Anything)
}

func (suite *ModelServiceTest) TestGetModelById_Success() {
	// given
	model := &models.Model{
		ID:          "1",
		Name:        "ModelName",
		Description: "ModelDescription",
		Dims:        []int{204, 1},
	}
	suite.modelRepo.On("FindByID", mock.Anything, "1").Return(model, nil)
	// when
	_, err := suite.modelService.GetModelById(context.Background(), "1")
	// then
	suite.NoError(err)
	suite.modelRepo.AssertCalled(suite.T(), "FindByID", mock.Anything, "1")
}

func (suite *ModelServiceTest) TestGetModelById_Error() {
	// given
	suite.modelRepo.On("FindByID", mock.Anything, "1").Return(nil, errors.New("error"))
	// when
	_, err := suite.modelService.GetModelById(context.Background(), "1")
	// then
	suite.Error(err)
	suite.modelRepo.AssertCalled(suite.T(), "FindByID", mock.Anything, "1")
}

func (suite *ModelServiceTest) TestGetModelByName_Success() {
	// given
	model := &models.Model{
		ID:          "1",
		Name:        "ModelName",
		Description: "ModelDescription",
		Dims:        []int{204, 1},
	}
	suite.modelRepo.On("FindByName", mock.Anything, "ModelName").Return(model, nil)
	// when
	_, err := suite.modelService.GetModelByName(context.Background(), "ModelName")
	// then
	suite.NoError(err)
	suite.modelRepo.AssertCalled(suite.T(), "FindByName", mock.Anything, "ModelName")
}

func (suite *ModelServiceTest) TestGetModelByNAme_Error() {
	// given
	suite.modelRepo.On("FindByName", mock.Anything, "ModelName").Return(nil, errors.New("error"))
	// when
	_, err := suite.modelService.GetModelByName(context.Background(), "ModelName")
	// then
	suite.Error(err)
	suite.modelRepo.AssertCalled(suite.T(), "FindByName", mock.Anything, "ModelName")
}

func (suite *ModelServiceTest) TestDeleteModelById_Success() {
	// given
	id := "someId"
	var transactionError error
	suite.modelRepo.On("InTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
		f := args.Get(1).(func(ctx context.Context) error)
		transactionError = f(context.Background())
	}).Return(nil)
	suite.modelRepo.On("DeleteByID", mock.Anything, id).Return(nil)
	suite.trainingRunService.On("DeleteAllByModelId", mock.Anything, id).Return(nil)
	// when
	suite.modelService.DeleteModelByID(context.Background(), id)
	// then
	suite.NoError(transactionError)
	suite.modelRepo.AssertCalled(suite.T(), "InTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error"))
	suite.modelRepo.AssertCalled(suite.T(), "DeleteByID", mock.Anything, id)
	suite.trainingRunService.AssertCalled(suite.T(), "DeleteAllByModelId", mock.Anything, id)
}

func (suite *ModelServiceTest) TestDeleteModelById_ModelError() {
	// given
	id := "someId"
	var transactionError error
	suite.modelRepo.On("InTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
		f := args.Get(1).(func(ctx context.Context) error)
		transactionError = f(context.Background())
	}).Return(nil)
	suite.modelRepo.On("DeleteByID", mock.Anything, id).Return(errors.New("error"))
	// when
	suite.modelService.DeleteModelByID(context.Background(), id)
	// then
	suite.Error(transactionError)
	suite.modelRepo.AssertCalled(suite.T(), "InTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error"))
	suite.modelRepo.AssertCalled(suite.T(), "DeleteByID", mock.Anything, id)
	suite.trainingRunService.AssertNotCalled(suite.T(), "DeleteAllByModelId", mock.Anything, id)
}

func (suite *ModelServiceTest) TestDeleteModelById_TrainingRunError() {
	// given
	id := "someId"
	var transactionError error
	suite.modelRepo.On("InTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
		f := args.Get(1).(func(ctx context.Context) error)
		transactionError = f(context.Background())
	}).Return(nil)
	suite.modelRepo.On("DeleteByID", mock.Anything, id).Return(nil)
	suite.trainingRunService.On("DeleteAllByModelId", mock.Anything, id).Return(errors.New("error"))
	// when
	suite.modelService.DeleteModelByID(context.Background(), id)
	// then
	suite.Error(transactionError)
	suite.modelRepo.AssertCalled(suite.T(), "InTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error"))
	suite.modelRepo.AssertCalled(suite.T(), "DeleteByID", mock.Anything, id)
	suite.trainingRunService.AssertCalled(suite.T(), "DeleteAllByModelId", mock.Anything, id)
}
