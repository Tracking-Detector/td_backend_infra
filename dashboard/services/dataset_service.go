package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Tracking-Detector/td_backend_infra/dashboard/config"
	"github.com/Tracking-Detector/td_backend_infra/dashboard/models"
)

type IDatasetService interface {
	GetAllDatasets() ([]*models.Dataset, error)
	GetDatasetByID(id string) (*models.Dataset, error)
	CreateDataset(datasetPayload *models.CreateDatasetPayload) (*models.Dataset, error)
	DeleteDataset(id string) error
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
	fmt.Println("Loading datasets")
	resp, err := s.restService.Get(s.dataSetServiceBaseUrl + "/datasets")
	if err != nil {
		s.loadingError = err
		return
	}
	apires := &models.APIResponse[[]*models.Dataset]{}
	if err := json.Unmarshal(resp.Body(), apires); err != nil {
		s.loadingError = errors.New(err.Error())
		return
	}
	if !apires.Success {
		s.loadingError = errors.New(apires.Message)
		return
	}
	s.cache = apires.Data
	fmt.Println("Found datasets", len(s.cache))
	s.loadingError = nil
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
	apires := &models.APIResponse[*models.Dataset]{}
	if err := json.Unmarshal(resp.Body(), apires); err != nil {
		return nil, errors.New(err.Error())
	}
	if !apires.Success {
		return nil, errors.New(apires.Message)
	}
	s.LoadAllDatasets()
	return apires.Data, nil
}

func (s *DatasetService) GetDatasetByID(id string) (*models.Dataset, error) {
	s.LoadAllDatasets()
	for _, dataset := range s.cache {
		if dataset.ID == id {
			return dataset, nil
		}
	}
	return nil, errors.New("Dataset not found")
}

func (s *DatasetService) DeleteDataset(id string) error {
	resp, err := s.restService.Delete(s.dataSetServiceBaseUrl + "/datasets/" + id)
	if err != nil {
		return err
	}
	apires := &models.APIResponse[interface{}]{}
	if err := json.Unmarshal(resp.Body(), apires); err != nil {
		return errors.New(err.Error())
	}
	if !apires.Success {
		return errors.New(apires.Message)
	}
	s.LoadAllDatasets()
	return nil
}
