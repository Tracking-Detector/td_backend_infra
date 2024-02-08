package service

import (
	"context"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
)

type IModelService interface {
	Save(ctx context.Context, model *models.Model) (*models.Model, error)
	GetAllModels(ctx context.Context) ([]*models.Model, error)
	GetModelByName(ctx context.Context, name string) (*models.Model, error)
	DeleteModelByID(ctx context.Context, id string) error
	GetModelById(ctx context.Context, id string) (*models.Model, error)
}

type ModelService struct {
	modelRepo          models.ModelRepository
	trainingrunService ITrainingrunService
}

func NewModelService(modelRepo models.ModelRepository, trainingrunService ITrainingrunService) *ModelService {
	return &ModelService{
		modelRepo:          modelRepo,
		trainingrunService: trainingrunService,
	}
}

func (s *ModelService) Save(ctx context.Context, model *models.Model) (*models.Model, error) {
	return s.modelRepo.Save(ctx, model)
}

func (s *ModelService) GetAllModels(ctx context.Context) ([]*models.Model, error) {
	return s.modelRepo.FindAll(ctx)
}

func (s *ModelService) GetModelByName(ctx context.Context, name string) (*models.Model, error) {
	return s.modelRepo.FindByName(ctx, name)
}

func (s *ModelService) GetModelById(ctx context.Context, id string) (*models.Model, error) {
	return s.modelRepo.FindByID(ctx, id)
}

func (s *ModelService) DeleteModelByID(ctx context.Context, id string) error {
	return s.modelRepo.InTransaction(ctx, func(ctx context.Context) error {
		if err := s.modelRepo.DeleteByID(ctx, id); err != nil {
			return err
		}
		return s.trainingrunService.DeleteAllByModelId(ctx, id)
	})
}
