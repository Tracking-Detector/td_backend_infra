package service

import (
	"context"
	"tds/shared/models"
)

type ITrainingrunService interface {
	FindAllTrainingRuns(ctx context.Context) ([]*models.TrainingRun, error)
	FindAllTrainingRunsForModelname(ctx context.Context, modelName string) ([]*models.TrainingRun, error)
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
