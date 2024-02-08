package acceptance

import (
	"encoding/json"
	"fmt"
	"net/http"

	"testing"
	"time"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/configs"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/controller"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/repository"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
	"github.com/Tracking-Detector/td_backend_infra/microservices/test/testsupport"
	"github.com/stretchr/testify/suite"
)

func TestDatasetControllerAcceptance(t *testing.T) {
	suite.Run(t, &DatasetControllerAcceptanceTest{})
}

type DatasetControllerAcceptanceTest struct {
	AcceptanceTest
	suite.Suite
	requestRepo       models.RequestRepository
	requestService    service.IRequestService
	datasetRepo       models.DatasetRepository
	datasetService    service.IDatasetService
	datasetController *controller.DatasetController
}

func (suite *DatasetControllerAcceptanceTest) SetupSuite() {
	suite.setupIntegration()
	suite.requestRepo = repository.NewMongoRequestRepository(configs.GetDatabase(configs.ConnectDB(suite.AcceptanceTest.ctx)))
	suite.requestService = service.NewRequestService(suite.requestRepo)
	suite.datasetRepo = repository.NewMongoDatasetRepository(configs.GetDatabase(configs.ConnectDB(suite.AcceptanceTest.ctx)))
	suite.datasetService = service.NewDatasetService(suite.datasetRepo, suite.requestRepo)
	suite.datasetController = controller.NewDatasetController(suite.datasetService)
	go func() {
		suite.datasetController.Start()
	}()
	time.Sleep(5 * time.Second)
}

func (suite *DatasetControllerAcceptanceTest) SetupTest() {
	suite.datasetRepo.DeleteAll(suite.AcceptanceTest.ctx)
}

func (suite *DatasetControllerAcceptanceTest) TearDownSuite() {
	suite.datasetController.Stop()
	suite.teardownIntegration()
}

func (suite *DatasetControllerAcceptanceTest) Test_FindByID() {
	// given
	dataset := &models.Dataset{
		Name:        "test",
		Description: "test",
		Label:       "test",
	}
	dataset, err := suite.datasetRepo.Save(suite.AcceptanceTest.ctx, dataset)
	// when
	newDataSet, err := suite.datasetRepo.FindByID(suite.AcceptanceTest.ctx, dataset.ID)
	// then
	fmt.Println(newDataSet, err)

}

func (suite *DatasetControllerAcceptanceTest) TestHealth_Success() {
	// given

	// when
	resp, err := testsupport.Get("http://localhost:8081/datasets/health")

	// then
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal(`{"success":true,"data":"System is running correct."}`, resp.Body)
}

func (suite *DatasetControllerAcceptanceTest) TestCreateDataset() {
	// given
	datasetPayload := &models.Dataset{
		Name:        "test",
		Description: "test",
		Label:       "test",
	}
	body, _ := json.Marshal(datasetPayload)
	// when
	resp, err := testsupport.Post("http://localhost:8081/datasets", string(body), "application/json")
	// then
	count, _ := suite.datasetRepo.Count(suite.AcceptanceTest.ctx)
	dataset, _ := suite.datasetRepo.FindByLabel(suite.AcceptanceTest.ctx, "test")
	suite.Equal(int64(1), count)
	suite.Equal("test", dataset.Label)
	suite.Equal("test", dataset.Name)
	suite.Equal("test", dataset.Description)
	suite.Equal(http.StatusCreated, resp.StatusCode)
	suite.NoError(err)
}

func (suite *DatasetControllerAcceptanceTest) TestGetAllDatasets() {
	// given
	datasetPayload := &models.Dataset{
		Name:        "test",
		Description: "test",
		Label:       "test",
	}
	body, _ := json.Marshal(datasetPayload)
	testsupport.Post("http://localhost:8081/datasets", string(body), "application/json")
	// when
	resp, err := testsupport.Get("http://localhost:8081/datasets")

	// then
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
}
