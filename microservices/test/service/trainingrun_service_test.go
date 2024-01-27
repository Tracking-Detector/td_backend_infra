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

func TestTrainingRunService(t *testing.T) {
	suite.Run(t, &TrainingRunServiceTest{})
}

type TrainingRunServiceTest struct {
	suite.Suite
	trainingRunService *service.TraingingrunService
	trainingrunRepo    *mocks.TrainingRunRepository
}

func (suite *TrainingRunServiceTest) SetupTest() {
	suite.trainingrunRepo = new(mocks.TrainingRunRepository)
	suite.trainingRunService = service.NewTraingingrunService(suite.trainingrunRepo)
}

func (suite *TrainingRunServiceTest) TestFindAllTrainingRuns_Success() {
	// given
	expectedResult := []*models.TrainingRun{{ID: "1", Name: "Run1"}, {ID: "2", Name: "Run2"}}
	suite.trainingrunRepo.On("FindAll", mock.Anything).Return(expectedResult, nil)

	// when
	result, err := suite.trainingRunService.FindAllTrainingRuns(context.Background())

	// then
	suite.trainingrunRepo.AssertCalled(suite.T(), "FindAll", mock.Anything)
	suite.NoError(err)
	suite.Equal(expectedResult, result)
}

func (suite *TrainingRunServiceTest) TestFindAllTrainingRuns_Error() {
	// given
	suite.trainingrunRepo.On("FindAll", mock.Anything).Return(nil, errors.New("error"))

	// when
	result, err := suite.trainingRunService.FindAllTrainingRuns(context.Background())

	// then
	suite.trainingrunRepo.AssertCalled(suite.T(), "FindAll", mock.Anything)
	suite.Error(err)
	suite.Nil(result)
}

func (suite *TrainingRunServiceTest) TestFindAllTrainingRunsForModelname_Success() {
	// given
	expectedResult := []*models.TrainingRun{{ID: "1", Name: "Run1"}, {ID: "2", Name: "Run2"}}
	modelName := "exampleModel"
	suite.trainingrunRepo.On("FindByModelName", mock.Anything, modelName).Return(expectedResult, nil)

	// when
	result, err := suite.trainingRunService.FindAllTrainingRunsForModelname(context.Background(), modelName)

	// then
	suite.trainingrunRepo.AssertCalled(suite.T(), "FindByModelName", mock.Anything, modelName)
	suite.NoError(err)
	suite.Equal(expectedResult, result)
}

func (suite *TrainingRunServiceTest) TestFindAllTrainingRunsForModelname_Error() {
	// given
	modelName := "exampleModel"
	suite.trainingrunRepo.On("FindByModelName", mock.Anything, modelName).Return(nil, errors.New("error"))

	// when
	result, err := suite.trainingRunService.FindAllTrainingRunsForModelname(context.Background(), modelName)

	// then
	suite.trainingrunRepo.AssertCalled(suite.T(), "FindByModelName", mock.Anything, modelName)
	suite.Error(err)
	suite.Nil(result)
}

func (suite *TrainingRunServiceTest) TestFindAllByModelId_Success() {
	// given
	expectedResult := []*models.TrainingRun{{ID: "1", Name: "Run1"}, {ID: "2", Name: "Run2"}}
	modelID := "exampleModelID"
	suite.trainingrunRepo.On("FindByModelID", mock.Anything, modelID).Return(expectedResult, nil)

	// when
	result, err := suite.trainingRunService.FindAllByModelId(context.Background(), modelID)

	// then
	suite.trainingrunRepo.AssertCalled(suite.T(), "FindByModelID", mock.Anything, modelID)
	suite.NoError(err)
	suite.Equal(expectedResult, result)
}

func (suite *TrainingRunServiceTest) TestFindAllByModelId_Error() {
	// given
	modelID := "exampleModelID"
	suite.trainingrunRepo.On("FindByModelID", mock.Anything, modelID).Return(nil, errors.New("error"))

	// when
	result, err := suite.trainingRunService.FindAllByModelId(context.Background(), modelID)

	// then
	suite.trainingrunRepo.AssertCalled(suite.T(), "FindByModelID", mock.Anything, modelID)
	suite.Error(err)
	suite.Nil(result)
}

func (suite *TrainingRunServiceTest) TestDeleteAllByModelId_Success() {
	// given
	modelID := "exampleModelID"
	suite.trainingrunRepo.On("DeleteAllByModelID", mock.Anything, modelID).Return(nil)

	// when
	err := suite.trainingRunService.DeleteAllByModelId(context.Background(), modelID)

	// then
	suite.trainingrunRepo.AssertCalled(suite.T(), "DeleteAllByModelID", mock.Anything, modelID)
	suite.NoError(err)
}

func (suite *TrainingRunServiceTest) TestDeleteAllByModelId_Error() {
	// given
	modelID := "exampleModelID"
	suite.trainingrunRepo.On("DeleteAllByModelID", mock.Anything, modelID).Return(errors.New("error"))

	// when
	err := suite.trainingRunService.DeleteAllByModelId(context.Background(), modelID)

	// then
	suite.trainingrunRepo.AssertCalled(suite.T(), "DeleteAllByModelID", mock.Anything, modelID)
	suite.Error(err)
}

func (suite *TrainingRunServiceTest) TestDeleteByID_Success() {
	// given
	runID := "exampleRunID"
	suite.trainingrunRepo.On("DeleteByID", mock.Anything, runID).Return(nil)

	// when
	err := suite.trainingRunService.DeleteByID(context.Background(), runID)

	// then
	suite.trainingrunRepo.AssertCalled(suite.T(), "DeleteByID", mock.Anything, runID)
	suite.NoError(err)
}

func (suite *TrainingRunServiceTest) TestDeleteByID_Error() {
	// given
	runID := "exampleRunID"
	suite.trainingrunRepo.On("DeleteByID", mock.Anything, runID).Return(errors.New("error"))

	// when
	err := suite.trainingRunService.DeleteByID(context.Background(), runID)

	// then
	suite.trainingrunRepo.AssertCalled(suite.T(), "DeleteByID", mock.Anything, runID)
	suite.Error(err)
}
