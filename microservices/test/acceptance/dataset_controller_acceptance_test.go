package acceptance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"tds/shared/configs"
	"tds/shared/controller"
	"tds/shared/models"
	"tds/shared/repository"
	"tds/shared/service"
	"tds/test/testsupport"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestDatasetControllerAcceptance(t *testing.T) {
	suite.Run(t, &DatasetControllerAcceptanceTest{})
}

type DatasetControllerAcceptanceTest struct {
	suite.Suite
	requestRepo       models.RequestRepository
	requestService    service.IRequestService
	datasetRepo       models.DatasetRepository
	datasetService    service.IDatasetService
	datasetController *controller.DatasetController
	ctx               context.Context
}

func (suite *DatasetControllerAcceptanceTest) SetupTest() {
	suite.ctx = context.Background()
	suite.requestRepo = repository.NewMongoRequestRepository(configs.GetDatabase(configs.ConnectDB(suite.ctx)))
	suite.requestService = service.NewRequestService(suite.requestRepo)
	suite.datasetRepo = repository.NewMongoDatasetRepository(configs.GetDatabase(configs.ConnectDB(suite.ctx)))
	suite.datasetService = service.NewDatasetService(suite.datasetRepo)
	suite.datasetController = controller.NewDatasetController(suite.datasetService)
	go func() {
		suite.datasetController.Start()
	}()
	time.Sleep(5 * time.Second)
}

func (suite *DatasetControllerAcceptanceTest) TearDownTest() {
	suite.datasetController.Stop()
}

func (suite *DatasetControllerAcceptanceTest) TestHealth_Success() {
	// given

	// when
	resp, err := testsupport.Get("http://localhost:8081/datasets/health")

	// then
	fmt.Println(resp, err)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	fmt.Println(resp.Body, "{\"message\":\"System is running correct.\",\"status\":200}")
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
	fmt.Println(resp, err)
	count, _ := suite.datasetRepo.Count(suite.ctx)
	dataset, _ := suite.datasetRepo.FindByLabel(suite.ctx, "test")
	suite.Equal(int64(1), count)
	suite.Equal("test", dataset.Label)
	suite.Equal("test", dataset.Name)
	suite.Equal("test", dataset.Description)
	suite.Equal(http.StatusCreated, resp.StatusCode)
	suite.NoError(err)
}
