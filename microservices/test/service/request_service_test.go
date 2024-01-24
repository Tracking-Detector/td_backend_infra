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

func TestRequestService(t *testing.T) {
	suite.Run(t, &RequestServiceTest{})
}

type RequestServiceTest struct {
	suite.Suite
	requestService *service.RequestService
	requestRepo    *mocks.RequestRepository
}

func (suite *RequestServiceTest) SetupTest() {
	suite.requestRepo = new(mocks.RequestRepository)
	suite.requestService = service.NewRequestService(suite.requestRepo)
}

func (suite *RequestServiceTest) TestGetRequestById_Success() {
	// given
	id := "someID"
	request := &models.RequestData{
		ID:         id,
		DocumentId: "documentId",
	}
	suite.requestRepo.On("FindByID", mock.Anything, id).Return(request, nil)
	// when
	_, err := suite.requestService.GetRequestById(context.Background(), id)
	// then
	suite.NoError(err)
	suite.requestRepo.AssertCalled(suite.T(), "FindByID", mock.Anything, id)
}

func (suite *RequestServiceTest) TestGetRequestById_Error() {
	// given
	id := "someID"

	suite.requestRepo.On("FindByID", mock.Anything, id).Return(nil, errors.New("error"))
	// when
	_, err := suite.requestService.GetRequestById(context.Background(), id)
	// then
	suite.Error(err)
	suite.requestRepo.AssertCalled(suite.T(), "FindByID", mock.Anything, id)
}

func (suite *RequestServiceTest) TestInsertManyRequest_Success() {
	// given
	requests := []*models.RequestData{{
		DocumentId: "documentId1",
	}, {
		DocumentId: "documentId2",
	},
	}
	suite.requestRepo.On("InsertMany", mock.Anything, requests).Return(nil)
	// when
	err := suite.requestService.InsertManyRequests(context.Background(), requests)
	// then
	suite.NoError(err)
	suite.requestRepo.AssertCalled(suite.T(), "InsertMany", mock.Anything, requests)
}

func (suite *RequestServiceTest) TestInsertManyRequest_Error() {
	// given
	requests := []*models.RequestData{{
		DocumentId: "documentId1",
	}, {
		DocumentId: "documentId2",
	},
	}
	suite.requestRepo.On("InsertMany", mock.Anything, requests).Return(errors.New("error"))
	// when
	err := suite.requestService.InsertManyRequests(context.Background(), requests)
	// then
	suite.Error(err)
	suite.requestRepo.AssertCalled(suite.T(), "InsertMany", mock.Anything, requests)
}

func (suite *RequestServiceTest) TestSaveRequest_Success() {
	// given
	id := "someID"
	request := &models.RequestData{
		ID:         id,
		DocumentId: "documentId",
	}
	suite.requestRepo.On("Save", mock.Anything, request).Return(nil)
	// when
	err := suite.requestService.SaveRequest(context.Background(), request)
	// then
	suite.NoError(err)
	suite.requestRepo.AssertCalled(suite.T(), "Save", mock.Anything, request)
}

func (suite *RequestServiceTest) TestSaveRequest_Error() {
	// given
	id := "someID"
	request := &models.RequestData{
		ID:         id,
		DocumentId: "documentId",
	}
	suite.requestRepo.On("Save", mock.Anything, request).Return(errors.New("error"))
	// when
	err := suite.requestService.SaveRequest(context.Background(), request)
	// then
	suite.Error(err)
	suite.requestRepo.AssertCalled(suite.T(), "Save", mock.Anything, request)
}

func (suite *RequestServiceTest) TestGetPagedRequestsFilterdByUrl_Success() {
	// given
	request := &models.RequestData{
		ID:         "someId",
		DocumentId: "documentId",
	}
	url := "someUrl/test"
	page := 10
	pageSize := 1
	suite.requestRepo.On("FindAllFilteredByUrlPaged", mock.Anything, url, page, pageSize).Return([]*models.RequestData{request}, nil)
	// when
	_, err := suite.requestService.GetPagedRequestsFilterdByUrl(context.Background(), url, page, pageSize)
	// then
	suite.NoError(err)
	suite.requestRepo.AssertCalled(suite.T(), "FindAllFilteredByUrlPaged", mock.Anything, url, page, pageSize)
}

func (suite *RequestServiceTest) TestGetPagedRequestsFilterdByUrl_Error() {
	// given

	url := "someUrl/test"
	page := 10
	pageSize := 1
	suite.requestRepo.On("FindAllFilteredByUrlPaged", mock.Anything, url, page, pageSize).Return(nil, errors.New("error"))
	// when
	_, err := suite.requestService.GetPagedRequestsFilterdByUrl(context.Background(), url, page, pageSize)
	// then
	suite.Error(err)
	suite.requestRepo.AssertCalled(suite.T(), "FindAllFilteredByUrlPaged", mock.Anything, url, page, pageSize)
}

func (suite *RequestServiceTest) TestCountDocumentsForUrlFilter_Success() {
	// given

	url := "someUrl/test"

	suite.requestRepo.On("CountDocuments", mock.Anything, url).Return(int64(210000000), nil)
	// when
	count, err := suite.requestService.CountDocumentsForUrlFilter(context.Background(), url)
	// then
	suite.Equal(int64(210000000), count)
	suite.NoError(err)
	suite.requestRepo.AssertCalled(suite.T(), "CountDocuments", mock.Anything, url)
}

func (suite *RequestServiceTest) TestCountDocumentsForUrlFilter_Error() {
	// given

	url := "someUrl/test"

	suite.requestRepo.On("CountDocuments", mock.Anything, url).Return(int64(0), errors.New("error"))
	// when
	count, err := suite.requestService.CountDocumentsForUrlFilter(context.Background(), url)
	// then
	suite.Equal(int64(0), count)
	suite.Error(err)
	suite.requestRepo.AssertCalled(suite.T(), "CountDocuments", mock.Anything, url)
}
