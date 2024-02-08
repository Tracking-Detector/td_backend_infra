package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Tracking-Detector/td_backend_infra/dashboard/config"
	"github.com/Tracking-Detector/td_backend_infra/dashboard/models"
)

type IDatasetService interface {
	GetAllDatasets() ([]*models.Dataset, error)
	CreateDataset(datasetPayload *models.CreateDatasetPayload) (*models.Dataset, error)
}

type DatasetService struct {
	restService           IRestService
	dataSetServiceBaseUrl string
	cache                 []*models.Dataset
	loadingError          error
	lastUpdate            time.Time
}

func NewDatasetService(restService IRestService) *DatasetService {
	return &DatasetService{
		restService:           restService,
		dataSetServiceBaseUrl: config.EnvDatasetServiceDomain(),
	}
}

func (s *DatasetService) LoadAllDatasets() {
	resp, err := s.restService.Get(s.dataSetServiceBaseUrl + "/datasets")
	if err != nil {
		s.loadingError = err
		return
	}
	result := resp.Result()
	if result == nil {
		s.cache = make([]*models.Dataset, 0)
		s.loadingError = nil
		return
	}
	apires := result.(*models.APIResponse)
	if !apires.Success {
		s.loadingError = errors.New(apires.Message)
		return
	}
	if datasets, ok := apires.Data.([]*models.Dataset); ok {
		s.cache = datasets
		s.loadingError = nil
	} else {
		s.loadingError = errors.New("Error parsing response")
	}
	s.lastUpdate = time.Now()
}

func (s *DatasetService) GetAllDatasets() ([]*models.Dataset, error) {
	if time.Since(s.lastUpdate) > 5*time.Minute || len(s.cache) == 0 || s.loadingError != nil {
		s.LoadAllDatasets()
	}
	return s.cache, s.loadingError
}

func (s *DatasetService) CreateDataset(datasetPayload *models.CreateDatasetPayload) (*models.Dataset, error) {
	resp, err := s.restService.Post(s.dataSetServiceBaseUrl+"/datasets", datasetPayload)
	fmt.Println(resp, err)
	if err != nil {
		return nil, err
	}
	apires := resp.Result().(*models.APIResponse)
	if !apires.Success {
		return nil, errors.New(apires.Message)
	}
	return apires.Data.(*models.Dataset), nil
}
