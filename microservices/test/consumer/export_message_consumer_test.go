package consumer_test

import (
	"os"
	"tds/shared/configs"
	"tds/shared/consumer"
	"tds/shared/messages"
	"tds/shared/models"
	"tds/test/testsupport/mocks"
	"testing"
	"time"

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
}

func (suite *ExportConsumerTest) SetupTest() {
	suite.internalJob = new(mocks.IExportJob)
	suite.externalJob = new(mocks.IExportJob)
	suite.queueAdapter = new(mocks.IQueueChannelAdapter)
	suite.exporterService = new(mocks.IExporterService)
	suite.exportRunService = new(mocks.IExportRunService)
	suite.exportConsumer = consumer.NewExportMessageConsumer(suite.internalJob, suite.externalJob, suite.exportRunService, suite.queueAdapter, suite.exporterService)
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
	jobs := []*messages.JobPayload{messages.NewJob("export", []string{"someId", "or", "dataset"})}
	suite.exportRunService.On("Save", mock.Anything, mock.Anything).Return(&models.ExportRun{
		ID:         "someRunId",
		ExporterId: exporter.ID,
		Name:       exporter.Name,
		Dataset:    "dataset",
		Start:      time.Now(),
		End:        time.Now(),
	}, nil)
	suite.queueAdapter.On("Consume", configs.EnvExportQueueName(), "", true, false, false, false, mock.Anything).Return(suite.createChan(jobs), nil)
	suite.exporterService.On("FindByID", mock.Anything, "someId").Return(exporter, nil)
	suite.internalJob.On("Execute", exporter, "or", "dataset").Return(nil)
	// when
	suite.exportConsumer.Consume()
	suite.exportConsumer.Wg.Wait()
	// then
	suite.queueAdapter.AssertCalled(suite.T(), "Consume", configs.EnvExportQueueName(), "", true, false, false, false, mock.Anything)
	suite.exporterService.AssertCalled(suite.T(), "FindByID", mock.Anything, "someId")
	suite.exportRunService.AssertNumberOfCalls(suite.T(), "Save", 2)
	suite.internalJob.AssertCalled(suite.T(), "Execute", exporter, "or", "dataset")
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
	jobs := []*messages.JobPayload{messages.NewJob("export", []string{"someId", "or", "dataset"})}
	suite.queueAdapter.On("Consume", configs.EnvExportQueueName(), "", true, false, false, false, mock.Anything).Return(suite.createChan(jobs), nil)
	suite.exportRunService.On("Save", mock.Anything, mock.Anything).Return(&models.ExportRun{
		ID:         "someRunId",
		ExporterId: exporter.ID,
		Name:       exporter.Name,
		Dataset:    "dataset",
		Start:      time.Now(),
		End:        time.Now(),
	}, nil)
	suite.exporterService.On("FindByID", mock.Anything, "someId").Return(exporter, nil)
	suite.externalJob.On("Execute", exporter, "or", "dataset").Return(nil)
	// when
	suite.exportConsumer.Consume()
	suite.exportConsumer.Wg.Wait()
	// then
	suite.queueAdapter.AssertCalled(suite.T(), "Consume", configs.EnvExportQueueName(), "", true, false, false, false, mock.Anything)
	suite.exporterService.AssertCalled(suite.T(), "FindByID", mock.Anything, "someId")
	suite.exportRunService.AssertNumberOfCalls(suite.T(), "Save", 2)
	suite.externalJob.AssertCalled(suite.T(), "Execute", exporter, "or", "dataset")
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
	jobs := []*messages.JobPayload{messages.NewJob("export", []string{"someId1", "or", "dataset"}),
		messages.NewJob("export", []string{"someId2", "or", "dataset"})}
	suite.exportRunService.On("Save", mock.Anything, mock.Anything).Return(&models.ExportRun{
		ID:         "someRunId",
		ExporterId: exporter1.ID,
		Name:       exporter1.Name,
		Dataset:    "dataset",
		Start:      time.Now(),
		End:        time.Now(),
	}, nil)
	suite.queueAdapter.On("Consume", configs.EnvExportQueueName(), "", true, false, false, false, mock.Anything).Return(suite.createChan(jobs), nil)
	suite.exporterService.On("FindByID", mock.Anything, "someId1").Return(exporter1, nil)
	suite.exporterService.On("FindByID", mock.Anything, "someId2").Return(exporter2, nil)
	suite.internalJob.On("Execute", exporter1, "or", "dataset").Return(nil)
	suite.externalJob.On("Execute", exporter2, "or", "dataset").Return(nil)
	// when
	suite.exportConsumer.Consume()
	suite.exportConsumer.Wg.Wait()
	// then
	suite.queueAdapter.AssertCalled(suite.T(), "Consume", configs.EnvExportQueueName(), "", true, false, false, false, mock.Anything)
	suite.exporterService.AssertCalled(suite.T(), "FindByID", mock.Anything, "someId1")
	suite.exporterService.AssertCalled(suite.T(), "FindByID", mock.Anything, "someId2")
	suite.exportRunService.AssertNumberOfCalls(suite.T(), "Save", 4)
	suite.internalJob.AssertCalled(suite.T(), "Execute", exporter1, "or", "dataset")
	suite.externalJob.AssertCalled(suite.T(), "Execute", exporter2, "or", "dataset")
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
