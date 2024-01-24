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
	exportConsumer  *consumer.ExportMessageConsumer
	internalJob     *mocks.IExportJob
	externalJob     *mocks.IExportJob
	queueAdapter    *mocks.IQueueChannelAdapter
	exporterService *mocks.IExporterService
}

func (suite *ExportConsumerTest) SetupTest() {
	suite.internalJob = new(mocks.IExportJob)
	suite.externalJob = new(mocks.IExportJob)
	suite.queueAdapter = new(mocks.IQueueChannelAdapter)
	suite.exporterService = new(mocks.IExporterService)
	suite.exportConsumer = consumer.NewExportMessageConsumer(suite.internalJob, suite.externalJob, suite.queueAdapter, suite.exporterService)
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
	suite.queueAdapter.On("Consume", configs.EnvExportQueueName(), "", true, false, false, false, mock.Anything).Return(suite.createChan(jobs), nil)
	suite.exporterService.On("FindByID", mock.Anything, "someId").Return(exporter, nil)
	suite.internalJob.On("Execute", exporter, "or", "dataset").Return(nil)
	// when
	suite.exportConsumer.Consume()
	// then
	time.Sleep(1 * time.Second)
	suite.queueAdapter.AssertCalled(suite.T(), "Consume", configs.EnvExportQueueName(), "", true, false, false, false, mock.Anything)
	suite.exporterService.AssertCalled(suite.T(), "FindByID", mock.Anything, "someId")
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
	suite.exporterService.On("FindByID", mock.Anything, "someId").Return(exporter, nil)
	suite.externalJob.On("Execute", exporter, "or", "dataset").Return(nil)
	// when
	suite.exportConsumer.Consume()
	// then
	time.Sleep(1 * time.Second)
	suite.queueAdapter.AssertCalled(suite.T(), "Consume", configs.EnvExportQueueName(), "", true, false, false, false, mock.Anything)
	suite.exporterService.AssertCalled(suite.T(), "FindByID", mock.Anything, "someId")
	suite.externalJob.AssertCalled(suite.T(), "Execute", exporter, "or", "dataset")
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
