package service

import (
	"context"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
)

type IExportRunService interface {
	Save(ctx context.Context, exportRun *models.ExportRun) (*models.ExportRun, error)
	GetAll(ctx context.Context) ([]*models.ExportRun, error)
	GetByExporterID(ctx context.Context, exporterId string) ([]*models.ExportRun, error)
	ExistByExporterIDAndRecducer(ctx context.Context, exporterId, reducer string) (bool, error)
	GetByID(ctx context.Context, id string) (*models.ExportRun, error)
}

type ExportRunService struct {
	exportRunRepository models.ExportRunRepository
}

func NewExportRunService(exportRunRepository models.ExportRunRepository) *ExportRunService {
	return &ExportRunService{
		exportRunRepository: exportRunRepository,
	}
}

func (s *ExportRunService) Save(ctx context.Context, exportRun *models.ExportRun) (*models.ExportRun, error) {
	return s.exportRunRepository.Save(ctx, exportRun)
}

func (s *ExportRunService) GetAll(ctx context.Context) ([]*models.ExportRun, error) {
	return s.exportRunRepository.FindAll(ctx)
}

func (s *ExportRunService) GetByExporterID(ctx context.Context, exporterId string) ([]*models.ExportRun, error) {
	return s.exportRunRepository.FindByExporterID(ctx, exporterId)
}

func (s *ExportRunService) GetByID(ctx context.Context, id string) (*models.ExportRun, error) {
	return s.exportRunRepository.FindByID(ctx, id)
}

func (s *ExportRunService) ExistByExporterIDAndRecducer(ctx context.Context, exporterId, reducer string) (bool, error) {
	return s.exportRunRepository.ExistByExporterIDAndRecducer(ctx, exporterId, reducer)
}
