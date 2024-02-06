package job_test

import (
	"io"
	"os"

	"testing"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/job"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	"github.com/Tracking-Detector/td_backend_infra/microservices/test/testsupport"
	"github.com/Tracking-Detector/td_backend_infra/microservices/test/testsupport/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestInternalExportJob(t *testing.T) {
	suite.Run(t, &InternalExportJobTest{})
}

type InternalExportJobTest struct {
	suite.Suite
	internalJob    *job.InternalExportJob
	requestRepo    *mocks.RequestRepository
	storageService *mocks.IStorageService
}

func (suite *InternalExportJobTest) SetupTest() {
	suite.requestRepo = new(mocks.RequestRepository)
	suite.storageService = new(mocks.IStorageService)
	suite.internalJob = job.NewInternalExportJob(suite.requestRepo, suite.storageService)
}

func (suite *InternalExportJobTest) TestExecute_Success() {
	// given
	requests := testsupport.LoadRequestJson()
	suite.Len(requests, 10)
	suite.requestRepo.On("StreamByDataset", mock.Anything, "").Return(testsupport.CreateResultsChannel(requests), testsupport.CreateErrorChannel([]error{}))
	suite.storageService.On("PutObject", mock.Anything, "", "GoExtractor204_EasyPrivacy_.csv.gz", mock.Anything, int64(-1), "application/gzip").Run(func(args mock.Arguments) {
		fileName := args.Get(2).(string)
		pr := args.Get(3).(io.Reader)
		file, _ := os.Create(fileName)
		io.Copy(file, pr)
	}).Return(nil)
	exporter := &models.Exporter{
		ID:          "someId",
		Name:        "GoExtractor204",
		Description: "someDescription",
		Dimensions:  []int{204, 1},
		Type:        models.IN_SERVICE,
	}
	// when
	metrics := suite.internalJob.Execute(exporter, "EasyPrivacy", "")
	// then
	suite.Assertions.Equal(10, metrics.Total)
	suite.Assertions.Equal(9, metrics.NonTracker)
	suite.Assertions.Equal(1, metrics.Tracker)
	expectedCsv := testsupport.LoadFile("../resources/requests/expected_encoding.csv")
	actualCsv := testsupport.LoadGzFile("./GoExtractor204_EasyPrivacy_.csv.gz")
	suite.Assertions.Equal(expectedCsv, actualCsv)
	os.Remove("./GoExtractor204_EasyPrivacy_.csv.gz")
}
