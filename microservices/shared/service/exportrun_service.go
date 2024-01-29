package service

import (
	"context"
	"tds/shared/models"
)

type IExportRunService interface {
	Save(ctx context.Context, exportRun *models.ExportRun) (*models.ExportRun, error)
	GetAll(ctx context.Context) ([]*models.ExportRun, error)
	GetByExporterID(ctx context.Context, exporterId string) ([]*models.ExportRun, error)
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
