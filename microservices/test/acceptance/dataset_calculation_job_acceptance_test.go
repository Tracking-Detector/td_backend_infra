package acceptance

import (
	"context"
	"encoding/json"
	"fmt"
	"tds/shared/configs"
	"tds/shared/job"
	"tds/shared/models"
	"tds/shared/repository"
	"tds/shared/service"
	"tds/test/testsupport"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestDatasetCalculationJobAcceptance(t *testing.T) {
	suite.Run(t, &DatasetCalculationJobAcceptanceTest{})
}

type DatasetCalculationJobAcceptanceTest struct {
	suite.Suite
	ctx              context.Context
	requestRepo      models.RequestRepository
	requestService   service.IRequestService
	datasetRepo      models.DatasetRepository
	datasetService   service.IDatasetService
	datasetMetricJob job.DatasetMetricJob
}

func (suite *DatasetCalculationJobAcceptanceTest) SetupTest() {
	suite.ctx = context.Background()
	suite.requestRepo = repository.NewMongoRequestRepository(configs.GetDatabase(configs.ConnectDB(suite.ctx)))
	suite.datasetRepo = repository.NewMongoDatasetRepository(configs.GetDatabase(configs.ConnectDB(suite.ctx)))
	suite.requestService = service.NewRequestService(suite.requestRepo)
	suite.datasetService = service.NewDatasetService(suite.datasetRepo)
	suite.datasetMetricJob = *job.NewDatasetMetricJob(suite.datasetService, suite.requestService)
	suite.requestRepo.DeleteAll(suite.ctx)
	suite.datasetRepo.DeleteAll(suite.ctx)
}

func (suite *DatasetCalculationJobAcceptanceTest) TestExecute_Success() {
	// given
	suite.datasetRepo.Save(suite.ctx, &models.Dataset{
		Name:        "Training dataset",
		Description: "This is a training dataset.",
		Label:       "train",
	})
	requests := testsupport.LoadRequestJson()
	for _, request := range requests {
		request.Dataset = "train"
	}
	suite.requestRepo.SaveAll(suite.ctx, requests)
	suite.datasetService.ReloadCache(suite.ctx)
	// when
	suite.datasetMetricJob.Execute()
	// then
	datasets := suite.datasetService.GetAllDatasets()
	suite.Equal(1, len(datasets))
	suite.Equal(10, datasets[0].Metrics.Total)
	suite.Equal(10, datasets[0].Metrics.ReducerMetric[0].Total)
	suite.Equal(10, datasets[0].Metrics.ReducerMetric[0].NonTracker)
	suite.Equal(0, datasets[0].Metrics.ReducerMetric[0].Tracker)
	suite.Equal(10, datasets[0].Metrics.ReducerMetric[1].Total)
	suite.Equal(9, datasets[0].Metrics.ReducerMetric[1].NonTracker)
	suite.Equal(1, datasets[0].Metrics.ReducerMetric[1].Tracker)
}

func (suite *DatasetCalculationJobAcceptanceTest) TestExecuteMultiple_Success() {
	// given
	suite.datasetRepo.Save(suite.ctx, &models.Dataset{
		Name:        "Training dataset",
		Description: "This is a training dataset.",
		Label:       "train",
	})
	suite.datasetRepo.Save(suite.ctx, &models.Dataset{
		Name:        "Test dataset",
		Description: "This is a test dataset.",
		Label:       "test",
	})
	requests := testsupport.LoadRequestJson()
	for i, request := range requests {
		if i%3 == 0 {
			request.Dataset = "train"
		} else {
			request.Dataset = "test"
		}
	}
	suite.requestRepo.SaveAll(suite.ctx, requests)
	suite.datasetService.ReloadCache(suite.ctx)
	// when
	suite.datasetMetricJob.Execute()
	// then
	datasets := suite.datasetService.GetAllDatasets()
	suite.Equal(2, len(datasets))
	t, _ := json.Marshal(datasets)
	fmt.Println(string(t))
	suite.Equal(4, datasets[0].Metrics.Total)
	suite.Equal(4, datasets[0].Metrics.ReducerMetric[0].Total)
	suite.Equal(4, datasets[0].Metrics.ReducerMetric[0].NonTracker)
	suite.Equal(0, datasets[0].Metrics.ReducerMetric[0].Tracker)
	suite.Equal(4, datasets[0].Metrics.ReducerMetric[1].Total)
	suite.Equal(3, datasets[0].Metrics.ReducerMetric[1].NonTracker)
	suite.Equal(1, datasets[0].Metrics.ReducerMetric[1].Tracker)
	suite.Equal(6, datasets[1].Metrics.Total)
	suite.Equal(6, datasets[1].Metrics.ReducerMetric[0].Total)
	suite.Equal(6, datasets[1].Metrics.ReducerMetric[0].NonTracker)
	suite.Equal(0, datasets[1].Metrics.ReducerMetric[0].Tracker)
	suite.Equal(6, datasets[1].Metrics.ReducerMetric[1].Total)
	suite.Equal(6, datasets[1].Metrics.ReducerMetric[1].NonTracker)
	suite.Equal(0, datasets[1].Metrics.ReducerMetric[1].Tracker)
}
