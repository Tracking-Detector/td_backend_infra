package service

import (
	"context"
	"tds/shared/models"
)

type ITrainingrunService interface {
	FindAllTrainingRuns(ctx context.Context) ([]*models.TrainingRun, error)
	FindAllTrainingRunsForModelname(ctx context.Context, modelName string) ([]*models.TrainingRun, error)
	FindAllByModelId(ctx context.Context, modelId string) ([]*models.TrainingRun, error)
	DeleteAllByModelId(ctx context.Context, id string) error
	DeleteByID(ctx context.Context, id string) error
}

type TraingingrunService struct {
	trainingrunRepo models.TrainingRunRepository
}

func NewTraingingrunService(trainingrunRepo models.TrainingRunRepository) *TraingingrunService {
	return &TraingingrunService{
		trainingrunRepo: trainingrunRepo,
	}
}

func (s *TraingingrunService) FindAllTrainingRuns(ctx context.Context) ([]*models.TrainingRun, error) {
	return s.trainingrunRepo.FindAll(ctx)
}

func (s *TraingingrunService) FindAllTrainingRunsForModelname(ctx context.Context, modelName string) ([]*models.TrainingRun, error) {
	return s.trainingrunRepo.FindByModelName(ctx, modelName)
}

func (s *TraingingrunService) FindAllByModelId(ctx context.Context, modelId string) ([]*models.TrainingRun, error) {
	return s.trainingrunRepo.FindByModelID(ctx, modelId)
}

func (s *TraingingrunService) DeleteAllByModelId(ctx context.Context, id string) error {
	return s.trainingrunRepo.DeleteAllByModelID(ctx, id)
}

func (s *TraingingrunService) DeleteByID(ctx context.Context, id string) error {
	return s.trainingrunRepo.DeleteByID(ctx, id)
}
