package consumer_test

import (
	"os"

	"testing"
	"time"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/configs"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/consumer"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/messages"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	"github.com/Tracking-Detector/td_backend_infra/microservices/test/testsupport/mocks"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestExportConsumer(t *testing.T) {
	suite.Run(t, &ExportConsumerTest{})
}

type ExportConsumerTest struct {
	suite.Suite
	exportConsumer   *consumer.ExportMessageConsumer
	internalJob      *mocks.IExportJob
	externalJob      *mocks.IExportJob
	queueAdapter     *mocks.IQueueChannelAdapter
	exportRunService *mocks.IExportRunService
	exporterService  *mocks.IExporterService
	datasetService   *mocks.IDatasetService
}

func (suite *ExportConsumerTest) SetupTest() {
	suite.internalJob = new(mocks.IExportJob)
	suite.externalJob = new(mocks.IExportJob)
	suite.queueAdapter = new(mocks.IQueueChannelAdapter)
	suite.exporterService = new(mocks.IExporterService)
	suite.exportRunService = new(mocks.IExportRunService)
	suite.datasetService = new(mocks.IDatasetService)
	suite.exportConsumer = consumer.NewExportMessageConsumer(suite.internalJob, suite.externalJob, suite.exportRunService, suite.queueAdapter, suite.exporterService, suite.datasetService)
}

func (suite *ExportConsumerTest) TestConsume_SuccessInternal() {
	// given
	os.Setenv("EXPORT_QUEUE", "export")
	exporter := &models.Exporter{
		ID:          "someId",
		Name:        "someName",
		Description: "someDescription",
		Dimensions:  []int{204, 1},
		Type:        models.IN_SERVICE,
	}
	dataset := &models.Dataset{
		ID:    "someId",
		Name:  "someName",
		Label: "",
	}
	jobs := []*messages.JobPayload{messages.NewJob("export", []string{"someId", "or", "someId"})}
	suite.datasetService.On("GetDatasetByID", mock.Anything, "someId").Return(dataset, nil)
	suite.exportRunService.On("Save", mock.Anything, mock.Anything).Return(&models.ExportRun{
		ID:         "someRunId",
		ExporterId: exporter.ID,
		Name:       exporter.Name,
		Dataset:    "someId",
		Start:      time.Now(),
		End:        time.Now(),
	}, nil)
	suite.queueAdapter.On("Consume", configs.EnvExportQueueName(), "ExportConsumer", true, false, false, false, mock.Anything).Return(suite.createChan(jobs), nil)
	suite.exporterService.On("FindByID", mock.Anything, "someId").Return(exporter, nil)
	suite.internalJob.On("Execute", exporter, "or", "").Return(nil)
	// when
	suite.exportConsumer.Consume()
	suite.exportConsumer.Wg.Wait()
	// then
	suite.queueAdapter.AssertCalled(suite.T(), "Consume", configs.EnvExportQueueName(), "ExportConsumer", true, false, false, false, mock.Anything)
	suite.exporterService.AssertCalled(suite.T(), "FindByID", mock.Anything, "someId")
	suite.exportRunService.AssertNumberOfCalls(suite.T(), "Save", 2)
	suite.internalJob.AssertCalled(suite.T(), "Execute", exporter, "or", "")
}

func (suite *ExportConsumerTest) TestConsume_SuccessExternal() {
	// given
	os.Setenv("EXPORT_QUEUE", "export")
	exporter := &models.Exporter{
		ID:          "someId",
		Name:        "someName",
		Description: "someDescription",
		Dimensions:  []int{204, 1},
		Type:        models.JS,
	}
	dataset := &models.Dataset{
		ID:    "someId",
		Name:  "someName",
		Label: "",
	}
	jobs := []*messages.JobPayload{messages.NewJob("export", []string{"someId", "or", "someId"})}
	suite.datasetService.On("GetDatasetByID", mock.Anything, "someId").Return(dataset, nil)
	suite.queueAdapter.On("Consume", configs.EnvExportQueueName(), "ExportConsumer", true, false, false, false, mock.Anything).Return(suite.createChan(jobs), nil)
	suite.exportRunService.On("Save", mock.Anything, mock.Anything).Return(&models.ExportRun{
		ID:         "someRunId",
		ExporterId: exporter.ID,
		Name:       exporter.Name,
		Dataset:    "someId",
		Start:      time.Now(),
		End:        time.Now(),
	}, nil)
	suite.exporterService.On("FindByID", mock.Anything, "someId").Return(exporter, nil)
	suite.externalJob.On("Execute", exporter, "or", "").Return(nil)
	// when
	suite.exportConsumer.Consume()
	suite.exportConsumer.Wg.Wait()
	// then
	suite.queueAdapter.AssertCalled(suite.T(), "Consume", configs.EnvExportQueueName(), "ExportConsumer", true, false, false, false, mock.Anything)
	suite.exporterService.AssertCalled(suite.T(), "FindByID", mock.Anything, "someId")
	suite.exportRunService.AssertNumberOfCalls(suite.T(), "Save", 2)
	suite.externalJob.AssertCalled(suite.T(), "Execute", exporter, "or", "")
}

func (suite *ExportConsumerTest) TestConsume_SuccessMultiple() {
	// given
	os.Setenv("EXPORT_QUEUE", "export")
	exporter1 := &models.Exporter{
		ID:          "someId1",
		Name:        "someName",
		Description: "someDescription",
		Dimensions:  []int{204, 1},
		Type:        models.IN_SERVICE,
	}
	exporter2 := &models.Exporter{
		ID:          "someId2",
		Name:        "someName",
		Description: "someDescription",
		Dimensions:  []int{204, 1},
		Type:        models.JS,
	}
	dataset := &models.Dataset{
		ID:    "someId",
		Name:  "someName",
		Label: "",
	}
	jobs := []*messages.JobPayload{messages.NewJob("export", []string{"someId1", "or", "someId"}),
		messages.NewJob("export", []string{"someId2", "or", "someId"})}
	suite.datasetService.On("GetDatasetByID", mock.Anything, "someId").Return(dataset, nil)
	suite.exportRunService.On("Save", mock.Anything, mock.Anything).Return(&models.ExportRun{
		ID:         "someRunId",
		ExporterId: exporter1.ID,
		Name:       exporter1.Name,
		Dataset:    "someId",
		Start:      time.Now(),
		End:        time.Now(),
	}, nil)
	suite.queueAdapter.On("Consume", configs.EnvExportQueueName(), "ExportConsumer", true, false, false, false, mock.Anything).Return(suite.createChan(jobs), nil)
	suite.exporterService.On("FindByID", mock.Anything, "someId1").Return(exporter1, nil)
	suite.exporterService.On("FindByID", mock.Anything, "someId2").Return(exporter2, nil)
	suite.internalJob.On("Execute", exporter1, "or", "").Return(nil)
	suite.externalJob.On("Execute", exporter2, "or", "").Return(nil)
	// when
	suite.exportConsumer.Consume()
	suite.exportConsumer.Wg.Wait()
	// then
	suite.queueAdapter.AssertCalled(suite.T(), "Consume", configs.EnvExportQueueName(), "ExportConsumer", true, false, false, false, mock.Anything)
	suite.exporterService.AssertCalled(suite.T(), "FindByID", mock.Anything, "someId1")
	suite.exporterService.AssertCalled(suite.T(), "FindByID", mock.Anything, "someId2")
	suite.exportRunService.AssertNumberOfCalls(suite.T(), "Save", 4)
	suite.internalJob.AssertCalled(suite.T(), "Execute", exporter1, "or", "")
	suite.externalJob.AssertCalled(suite.T(), "Execute", exporter2, "or", "")
}

func (suite *ExportConsumerTest) createChan(jobs []*messages.JobPayload) <-chan amqp.Delivery {
	jobCh := make(chan amqp.Delivery, len(jobs))
	for _, job := range jobs {
		ser, _ := job.Serialize()
		jobCh <- amqp.Delivery{
			ContentType:  "text/plain",
			Body:         []byte(ser),
			DeliveryMode: amqp.Persistent,
		}
	}
	defer close(jobCh)
	return jobCh
}
