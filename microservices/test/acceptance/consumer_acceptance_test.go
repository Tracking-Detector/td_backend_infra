package acceptance

import (
	"os"

	"testing"
	"time"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/configs"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/consumer"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/job"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/queue"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/repository"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/storage"
	"github.com/Tracking-Detector/td_backend_infra/microservices/test/testsupport"
	"github.com/stretchr/testify/suite"
)

func TestConsumerAcceptance(t *testing.T) {
	suite.Run(t, &ExportConsumerAcceptanceTest{})
}

type ExportConsumerAcceptanceTest struct {
	AcceptanceTest
	suite.Suite
	storageService   *service.MinIOStorageService
	publisherService *service.PublishService
	datasetService   *service.DatasetService
	requestRepo      *repository.MongoRequestRepository
	datasetRepo      *repository.MongoDatasetRepository
	exporterRepo     *repository.MongoExporterRepository
	exportRunRepo    *repository.MongoExportRunRunRepository
	exportRunService *service.ExportRunService
	exportConsumer   *consumer.ExportMessageConsumer
	queueAdapter     *queue.RabbitMQChannelAdapter
}

func (suite *ExportConsumerAcceptanceTest) SetupSuite() {
	suite.setupIntegration()
	suite.exporterRepo = repository.NewMongoExporterRepository(configs.GetDatabase(configs.ConnectDB(suite.ctx)))
	suite.requestRepo = repository.NewMongoRequestRepository(configs.GetDatabase(configs.ConnectDB(suite.ctx)))
	suite.exportRunRepo = repository.NewMongoExportRunRunRepository(configs.GetDatabase(configs.ConnectDB(suite.ctx)))
	suite.datasetRepo = repository.NewMongoDatasetRepository(configs.GetDatabase(configs.ConnectDB(suite.ctx)))
	suite.datasetService = service.NewDatasetService(suite.datasetRepo, suite.requestRepo)
	suite.exportRunService = service.NewExportRunService(suite.exportRunRepo)
	minioClient := configs.ConnectMinio()
	rabbitMqChannel := configs.ConnectRabbitMQ()
	suite.queueAdapter = queue.NewRabbitMQChannelAdapter(rabbitMqChannel)
	minioStorageAdapter := storage.NewMinIOStorageAdapter(minioClient)
	suite.storageService = service.NewMinIOStorageService(minioStorageAdapter)
	suite.publisherService = service.NewPublishService(suite.queueAdapter)
	suite.datasetRepo.DeleteAll(suite.ctx)
	internalJob := job.NewInternalExportJob(suite.requestRepo, suite.storageService)
	externJob := job.NewExternalExportJob(suite.requestRepo, suite.storageService)
	suite.exportConsumer = consumer.NewExportMessageConsumer(internalJob, externJob,
		suite.exportRunService, suite.queueAdapter, service.NewExporterService(suite.exporterRepo), suite.datasetService)
	go func() {
		suite.exportConsumer.Consume()
	}()
}

func (suite *ExportConsumerAcceptanceTest) SetupTest() {
	suite.requestRepo.DeleteAll(suite.ctx)
	suite.exporterRepo.DeleteAll(suite.ctx)
	suite.exportRunRepo.DeleteAll(suite.ctx)
	suite.queueAdapter.PurgeQueue(configs.EnvExportBucketName(), false)
}

func (suite *ExportConsumerAcceptanceTest) TearDownSuite() {
	suite.exportConsumer.Stop()
	suite.teardownIntegration()
}

func (suite *ExportConsumerAcceptanceTest) TestExportConsumer_ForExternalExporterSuccess() {
	// given
	suite.storageService.VerifyBucketExists(suite.ctx, configs.EnvExportBucketName())
	suite.storageService.VerifyBucketExists(suite.ctx, configs.EnvExtractorBucketName())
	extractorFilePath := "../resources/exporter/exporter204.js"
	fileLoc := "exporter204.js"
	file, _ := os.Open(extractorFilePath)
	suite.storageService.PutObject(suite.ctx, configs.EnvExtractorBucketName(), "exporter204.js", file, -1, "application/javascript")
	exporter := &models.Exporter{
		ID:                   "someId",
		Name:                 "someName",
		Description:          "someDescription",
		Dimensions:           []int{204, 1},
		Type:                 models.JS,
		ExportScriptLocation: &fileLoc,
	}
	dataset, _ := suite.datasetRepo.Save(suite.ctx, &models.Dataset{
		Name:        "someName",
		Description: "someDescription",
		Label:       "",
	})
	suite.exporterRepo.Save(suite.ctx, exporter)
	requests := testsupport.LoadRequestJson()
	suite.requestRepo.SaveAll(suite.ctx, requests)

	// when
	suite.publisherService.EnqueueExportJob("someId", "EasyPrivacy", dataset.ID)
	time.Sleep(1 * time.Second)
	// then
	suite.exportConsumer.Wg.Wait()
	export, err := suite.storageService.GetObject(suite.ctx, configs.EnvExportBucketName(), "someName_EasyPrivacy_.csv.gz")
	suite.NoError(err)
	expectedCsv := testsupport.LoadFile("../resources/requests/expected_encoding.csv")
	actualCsv := testsupport.Unzip(export)
	suite.Equal(expectedCsv, actualCsv)
	count, _ := suite.exportRunRepo.Count(suite.ctx)
	suite.Equal(int64(1), count)
	exportRuns, _ := suite.exportRunRepo.FindAll(suite.ctx)
	suite.Equal("someId", exportRuns[0].ExporterId)
	suite.Equal("someName", exportRuns[0].Name)
	suite.Equal("EasyPrivacy", exportRuns[0].Reducer)
	suite.Equal(dataset.ID, exportRuns[0].Dataset)
	suite.Equal(9, exportRuns[0].Metrics.NonTracker)
	suite.Equal(1, exportRuns[0].Metrics.Tracker)
	suite.Equal(10, exportRuns[0].Metrics.Total)

}
