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

func TestExporterService(t *testing.T) {
	suite.Run(t, &ExporterServiceTest{})
}

type ExporterServiceTest struct {
	suite.Suite
	exporterService *service.ExporterService
	exporterRepo    *mocks.ExporterRepository
}

func (suite *ExporterServiceTest) SetupTest() {
	suite.exporterRepo = new(mocks.ExporterRepository)
	suite.exporterService = service.NewExporterService(suite.exporterRepo)
}

func (suite *ExporterServiceTest) TestGetAllExporter() {
	// given
	exporters := []*models.Exporter{{ID: "1", Name: "Exporter1"}, {ID: "2", Name: "Exporter2"}}
	suite.exporterRepo.On("FindAll", mock.Anything).Return(exporters, nil)

	// when
	result, err := suite.exporterService.GetAllExporter(context.Background())

	// then
	suite.exporterRepo.AssertCalled(suite.T(), "FindAll", mock.Anything)
	suite.NoError(err)
	suite.Equal(exporters, result)
}

func (suite *ExporterServiceTest) TestInitInCodeExports() {
	// given
	suite.exporterRepo.On("Save", mock.Anything, mock.AnythingOfType("*models.Exporter")).Return(&models.Exporter{}, nil)
	// when
	suite.exporterService.InitInCodeExports(context.Background())

	// then
	suite.exporterRepo.AssertCalled(suite.T(), "Save", mock.Anything, mock.AnythingOfType("*models.Exporter"))
}

func (suite *ExporterServiceTest) TestIsValidExporter_Valid() {
	// given
	suite.exporterRepo.On("FindByID", mock.Anything, "validExporter").Return(&models.Exporter{ID: "validExporter"}, nil)

	// when
	result := suite.exporterService.IsValidExporter(context.Background(), "validExporter")

	// then
	suite.exporterRepo.AssertCalled(suite.T(), "FindByID", mock.Anything, "validExporter")
	suite.True(result)
}

func (suite *ExporterServiceTest) TestIsValidExporter_Invalid() {
	// given
	suite.exporterRepo.On("FindByID", mock.Anything, "invalidExporter").Return(nil, errors.New("not found"))

	// when
	result := suite.exporterService.IsValidExporter(context.Background(), "invalidExporter")

	// then
	suite.exporterRepo.AssertCalled(suite.T(), "FindByID", mock.Anything, "invalidExporter")
	suite.False(result)
}

func (suite *ExporterServiceTest) TestFindByID() {
	// given
	expectedExporter := &models.Exporter{ID: "1", Name: "Exporter1"}
	suite.exporterRepo.On("FindByID", mock.Anything, "1").Return(expectedExporter, nil)

	// when
	result, err := suite.exporterService.FindByID(context.Background(), "1")

	// then
	suite.exporterRepo.AssertCalled(suite.T(), "FindByID", mock.Anything, "1")
	suite.NoError(err)
	suite.Equal(expectedExporter, result)
}
