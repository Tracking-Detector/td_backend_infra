package acceptance

import (
	"fmt"

	"testing"
	"time"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/configs"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/controller"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/queue"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/repository"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
	"github.com/Tracking-Detector/td_backend_infra/microservices/test/testsupport"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/suite"
)

func TestDispatchControllerAcceptance(t *testing.T) {
	suite.Run(t, &DispatchControllerAcceptanceTest{})
}

type DispatchControllerAcceptanceTest struct {
	suite.Suite
	AcceptanceTest
	publishController  *controller.DispatchController
	exporterService    service.IExporterService
	exportRunService   service.IExportRunService
	publishService     service.IPublishService
	exporterRepo       models.ExporterRepository
	requestRepo        models.RequestRepository
	exportRunRepo      models.ExportRunRepository
	trainingRunRepo    models.TrainingRunRepository
	modelRepo          models.ModelRepository
	datasetRepo        models.DatasetRepository
	datasetService     service.IDatasetService
	queueAdapter       queue.IQueueChannelAdapter
	modelService       service.IModelService
	trainingRunService service.ITrainingrunService
	testConsumer       *testsupport.TestQueueConsumer
	rabbitMq           *amqp.Channel
}

func (suite *DispatchControllerAcceptanceTest) SetupSuite() {
	suite.setupIntegration()
	mongoClient := configs.ConnectDB(suite.ctx)
	suite.rabbitMq = configs.ConnectRabbitMQ()
	suite.requestRepo = repository.NewMongoRequestRepository(configs.GetDatabase(mongoClient))
	suite.queueAdapter = queue.NewRabbitMQChannelAdapter(suite.rabbitMq)
	suite.exporterRepo = repository.NewMongoExporterRepository(configs.GetDatabase(mongoClient))
	suite.exportRunRepo = repository.NewMongoExportRunRunRepository(configs.GetDatabase(mongoClient))
	suite.modelRepo = repository.NewMongoModelRepository(configs.GetDatabase(mongoClient))
	suite.trainingRunRepo = repository.NewMongoTrainingRunRepository(configs.GetDatabase(mongoClient))
	suite.datasetRepo = repository.NewMongoDatasetRepository(configs.GetDatabase(mongoClient))

	suite.trainingRunService = service.NewTraingingrunService(suite.trainingRunRepo)
	suite.datasetService = service.NewDatasetService(suite.datasetRepo, suite.requestRepo)
	suite.exportRunService = service.NewExportRunService(suite.exportRunRepo)
	suite.exporterService = service.NewExporterService(suite.exporterRepo)
	suite.modelService = service.NewModelService(suite.modelRepo, suite.trainingRunService)
	suite.publishService = service.NewPublishService(suite.queueAdapter)
	suite.testConsumer = testsupport.NewTestQueueConsumer(queue.NewRabbitMQChannelAdapter(suite.rabbitMq))
	suite.publishController = controller.NewDispatchController(suite.exporterService,
		suite.publishService, suite.modelService, suite.datasetService, suite.exportRunService)
	go func() {
		suite.publishController.Start()
	}()
	time.Sleep(5 * time.Second)
}

func (suite *DispatchControllerAcceptanceTest) SetupTest() {
	suite.exporterRepo.DeleteAll(suite.ctx)
	suite.modelRepo.DeleteAll(suite.ctx)
	suite.trainingRunRepo.DeleteAll(suite.ctx)
	suite.datasetRepo.DeleteAll(suite.ctx)
	suite.exportRunRepo.DeleteAll(suite.ctx)
	suite.queueAdapter.PurgeQueue(configs.EnvExportQueueName(), false)
	suite.queueAdapter.PurgeQueue(configs.EnvTrainQueueName(), false)
}

func (suite *DispatchControllerAcceptanceTest) TearDownTest() {
	suite.testConsumer.Stop()

}

func (suite *DispatchControllerAcceptanceTest) TearDownSuite() {
	suite.publishController.Stop()
	suite.teardownIntegration()
}

func (suite *DispatchControllerAcceptanceTest) TestHealth_Success() {
	// given
	// when
	resp, err := testsupport.Get("http://localhost:8081/dispatch/health")

	// then
	suite.NoError(err)
	suite.Equal(200, resp.StatusCode)
}

func (suite *DispatchControllerAcceptanceTest) TestDispatchExportJob_Success() {
	// given
	loc := "exporter.js"
	exporter, _ := suite.exporterRepo.Save(suite.ctx, &models.Exporter{
		Name:                 "ExporterJs204",
		Description:          "ExporterJs204",
		Dimensions:           []int{204, 1},
		Type:                 models.JS,
		ExportScriptLocation: &loc,
	})
	dataset, _ := suite.datasetRepo.Save(suite.ctx, &models.Dataset{
		Name:        "Verification",
		Description: "Can be used for verifaction",
		Label:       "verifiaction",
	})
	reducer := "or"
	// when
	go suite.testConsumer.Consume(configs.EnvExportQueueName(), 1)
	time.Sleep(5 * time.Second)
	resp, err := testsupport.Post(fmt.Sprintf("http://localhost:8081/dispatch/export/%s/%s/%s", exporter.ID, reducer, dataset.ID), "", "")
	suite.testConsumer.WaitForMessages(configs.EnvExportQueueName(), 1)
	// then
	suite.NoError(err)
	suite.Equal(201, resp.StatusCode)
	suite.Equal(1, len(suite.testConsumer.QueueMessages[configs.EnvExportQueueName()]))
	suite.Equal(fmt.Sprintf(`{"functionName":"export","args":["%s","%s","%s"]}`, exporter.ID, reducer, dataset.ID), suite.testConsumer.QueueMessages[configs.EnvExportQueueName()][0])
}

// func (suite *DispatchControllerAcceptanceTest) TestDispatchExportJob_ErrorReducerNotFound() {
// 	// given
// 	loc := "exporter.js"
// 	exporter, _ := suite.exporterRepo.Save(suite.ctx, &models.Exporter{
// 		Name:                 "ExporterJs204",
// 		Description:          "ExporterJs204",
// 		Dimensions:           []int{204, 1},
// 		Type:                 models.JS,
// 		ExportScriptLocation: &loc,
// 	})
// 	dataset, _ := suite.datasetRepo.Save(suite.ctx, &models.Dataset{
// 		Name:        "Verification",
// 		Description: "Can be used for verifaction",
// 		Label:       "verifiaction",
// 	})
// 	reducer := "notKnown"
// 	// when
// 	resp, err := testsupport.Post(fmt.Sprintf("http://localhost:8081/dispatch/export/%s/%s/%s", exporter.ID, reducer, dataset.ID), "", "")
// 	// then
// 	suite.NoError(err)
// 	suite.Equal(400, resp.StatusCode)
// 	suite.Equal(`{"success":false,"message":"The reducer type is not valid"}`, resp.Body)
// }

// func (suite *DispatchControllerAcceptanceTest) TestDispatchExportJob_ErrorDatasetNotFound() {
// 	loc := "exporter.js"
// 	exporter, _ := suite.exporterRepo.Save(suite.ctx, &models.Exporter{
// 		Name:                 "ExporterJs204",
// 		Description:          "ExporterJs204",
// 		Dimensions:           []int{204, 1},
// 		Type:                 models.JS,
// 		ExportScriptLocation: &loc,
// 	})
// 	datasetNotInDbId := "5f5e7e3e3e3e3e3e3e3e3e3e"
// 	reducer := "or"
// 	// when
// 	resp, err := testsupport.Post(fmt.Sprintf("http://localhost:8081/dispatch/export/%s/%s/%s", exporter.ID, reducer, datasetNotInDbId), "", "")
// 	// then
// 	suite.NoError(err)
// 	suite.Equal(404, resp.StatusCode)
// 	suite.Equal(`{"success":false,"message":"The dataset for the given id does not exist."}`, resp.Body)
// }

// func (suite *DispatchControllerAcceptanceTest) TestDispatchExportJob_ErrorExporterNotFound() {
// 	// given
// 	randomExporterId := "5f5e7e3e3e3e3e3e3e3e3e3e"
// 	dataset, _ := suite.datasetRepo.Save(suite.ctx, &models.Dataset{
// 		Name:        "Verification",
// 		Description: "Can be used for verifaction",
// 		Label:       "verifiaction",
// 	})
// 	reducer := "or"
// 	// when
// 	resp, err := testsupport.Post(fmt.Sprintf("http://localhost:8081/dispatch/export/%s/%s/%s", randomExporterId, reducer, dataset.ID), "", "")
// 	// then
// 	suite.NoError(err)
// 	suite.Equal(404, resp.StatusCode)
// 	suite.Equal(`{"success":false,"message":"The extractor for the given id does not exist."}`, resp.Body)
// }

func (suite *DispatchControllerAcceptanceTest) TestDispatchTrainingJob_Success() {
	// given
	model, _ := suite.modelRepo.Save(suite.ctx, &models.Model{
		Name:        "Model1",
		Description: "Model1",
		Dims:        []int{204, 1},
	})
	exporter, _ := suite.exporterRepo.Save(suite.ctx, &models.Exporter{
		Name:                 "ExporterJs204",
		Description:          "ExporterJs204",
		Dimensions:           []int{204, 1},
		Type:                 models.JS,
		ExportScriptLocation: nil,
	})
	reducer := "or"
	suite.exportRunRepo.Save(suite.ctx, &models.ExportRun{
		ExporterId: exporter.ID,
		Reducer:    reducer,
		Start:      time.Now(),
		End:        time.Now(),
	})
	// when
	go suite.testConsumer.Consume(configs.EnvTrainQueueName(), 1)
	time.Sleep(5 * time.Second)
	resp, err := testsupport.Post(fmt.Sprintf("http://localhost:8081/dispatch/train/%s/run/%s/%s", model.ID, exporter.ID, reducer), "", "")
	suite.testConsumer.WaitForMessages(configs.EnvTrainQueueName(), 1)
	// then
	suite.NoError(err)
	suite.Equal(201, resp.StatusCode)
	// suite.Equal(1, len(testConsumer.QueueMessages[configs.EnvTrainQueueName()]))
	// suite.Equal(fmt.Sprintf(`{"functionName":"train_model","args":["%s","%s","%s"]}`, model.ID, exporter.ID, reducer), testConsumer.QueueMessages[configs.EnvTrainQueueName()][0])
}

// func (suite *DispatchControllerAcceptanceTest) TestDispatchTrainingJob_ErrorNoRunFound() {
// 	// given
// 	model, _ := suite.modelRepo.Save(suite.ctx, &models.Model{
// 		Name:        "Model1",
// 		Description: "Model1",
// 		Dims:        []int{204, 1},
// 	})
// 	exporter, _ := suite.exporterRepo.Save(suite.ctx, &models.Exporter{
// 		Name:                 "ExporterJs204",
// 		Description:          "ExporterJs204",
// 		Dimensions:           []int{204, 1},
// 		Type:                 models.JS,
// 		ExportScriptLocation: nil,
// 	})
// 	reducer := "or"
// 	// when
// 	resp, err := testsupport.Post(fmt.Sprintf("http://localhost:8081/dispatch/train/%s/run/%s/%s", model.ID, exporter.ID, reducer), "", "")

// 	// then
// 	suite.NoError(err)
// 	suite.Equal(404, resp.StatusCode)
// 	suite.Equal(`{"success":false,"message":"The export for the given id and reducer does not exist."}`, resp.Body)
// }

// func (suite *DispatchControllerAcceptanceTest) TestDispatchTrainingJob_ErrorModelNotFound() {
// 	// given
// 	random := "5f5e7e3e3e3e3e3e3e3e3e3e"
// 	exporter, _ := suite.exporterRepo.Save(suite.ctx, &models.Exporter{
// 		Name:                 "ExporterJs204",
// 		Description:          "ExporterJs204",
// 		Dimensions:           []int{204, 1},
// 		Type:                 models.JS,
// 		ExportScriptLocation: nil,
// 	})
// 	reducer := "or"
// 	suite.exportRunRepo.Save(suite.ctx, &models.ExportRun{
// 		ExporterId: exporter.ID,
// 		Reducer:    reducer,
// 		Start:      time.Now(),
// 		End:        time.Now(),
// 	})
// 	// when
// 	resp, err := testsupport.Post(fmt.Sprintf("http://localhost:8081/dispatch/train/%s/run/%s/%s", random, exporter.ID, reducer), "", "")
// 	// then
// 	suite.NoError(err)
// 	suite.Equal(404, resp.StatusCode)
// 	suite.Equal(`{"success":false,"message":"The model for the given id does not exist."}`, resp.Body)
// }

// func (suite *DispatchControllerAcceptanceTest) TestDispatchTrainingJob_ErrorDimensionMismatch() {
// 	// given
// 	model, _ := suite.modelRepo.Save(suite.ctx, &models.Model{
// 		Name:        "Model1",
// 		Description: "Model1",
// 		Dims:        []int{204, 1},
// 	})
// 	exporter, _ := suite.exporterRepo.Save(suite.ctx, &models.Exporter{
// 		Name:                 "ExporterJs204",
// 		Description:          "ExporterJs204",
// 		Dimensions:           []int{204, 2},
// 		Type:                 models.JS,
// 		ExportScriptLocation: nil,
// 	})
// 	reducer := "or"
// 	suite.exportRunRepo.Save(suite.ctx, &models.ExportRun{
// 		ExporterId: exporter.ID,
// 		Reducer:    reducer,
// 		Start:      time.Now(),
// 		End:        time.Now(),
// 	})
// 	// when
// 	resp, err := testsupport.Post(fmt.Sprintf("http://localhost:8081/dispatch/train/%s/run/%s/%s", model.ID, exporter.ID, reducer), "", "")
// 	// then
// 	suite.NoError(err)
// 	suite.Equal(400, resp.StatusCode)
// 	suite.Equal(`{"success":false,"message":"There is a dimension mismatch for the dataset and the model."}`, resp.Body)
// }

// func (suite *DispatchControllerAcceptanceTest) TestDispatchTrainingJob_ErrorExporterNotFound() {
// 	// given
// 	random := "5f5e7e3e3e3e3e3e3e3e3e3e"
// 	model, _ := suite.modelRepo.Save(suite.ctx, &models.Model{
// 		Name:        "Model1",
// 		Description: "Model1",
// 		Dims:        []int{204, 1},
// 	})
// 	reducer := "or"
// 	// when
// 	resp, err := testsupport.Post(fmt.Sprintf("http://localhost:8081/dispatch/train/%s/run/%s/%s", model.ID, random, reducer), "", "")
// 	// then
// 	suite.NoError(err)
// 	suite.Equal(404, resp.StatusCode)
// 	suite.Equal(`{"success":false,"message":"The extractor for the given id does not exist."}`, resp.Body)
// }
