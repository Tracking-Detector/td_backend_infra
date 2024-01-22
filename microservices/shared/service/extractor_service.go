package service

import (
	"context"
	"tds/shared/extractor"
	"tds/shared/models"
)

type IExporterService interface {
	GetAllExporter(ctx context.Context) ([]*models.Exporter, error)
	InitInCodeExports(ctx context.Context)
}

type ExporterService struct {
	extractorRepo models.ExporterRepository
}

func NewExporterService(extractorRepo models.ExporterRepository) *ExporterService {
	return &ExporterService{
		extractorRepo: extractorRepo,
	}
}

func (s *ExporterService) GetAllExporter(ctx context.Context) ([]*models.Exporter, error) {
	return s.extractorRepo.FindAll(ctx)
}

func (s *ExporterService) InitInCodeExports(ctx context.Context) {
	for _, ext := range extractor.EXTRACTORS {
		exporterData := models.Exporter{
			Name:        ext.GetName(),
			Description: ext.GetDescription(),
			Dimensions:  ext.GetDimensions(),
			Type:        models.IN_SERVICE,
		}
		s.extractorRepo.Save(ctx, &exporterData)
	}
}
