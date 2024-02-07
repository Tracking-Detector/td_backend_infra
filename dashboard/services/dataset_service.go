package services

import (
	"time"

	"github.com/Tracking-Detector/td_backend_infra/dashboard/config"
	"github.com/Tracking-Detector/td_backend_infra/dashboard/models"
)

type IDatasetService interface {
}

type DatasetService struct {
	restService           IRestService
	dataSetServiceBaseUrl string
	cache                 []*models.ServiceStatus
	lastUpdate            time.Time
}

func NewDatasetService(restService IRestService) *DatasetService {
	return &DatasetService{
		restService:           restService,
		dataSetServiceBaseUrl: config.EnvDatasetServiceDomain(),
	}
}
