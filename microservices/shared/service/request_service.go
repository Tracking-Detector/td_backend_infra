package service

import (
	"context"
	"tds/shared/models"
)

type IRequestService interface {
	GetRequestById(ctx context.Context, id string) (*models.RequestData, error)
	InsertManyRequests(ctx context.Context, requests []*models.RequestData) error
	SaveRequest(ctx context.Context, request *models.RequestData) (*models.RequestData, error)
	GetPagedRequestsFilterdByUrl(ctx context.Context, url string, page, pageSize int) ([]*models.RequestData, error)
	CountDocumentsForUrlFilter(ctx context.Context, url string) (int64, error)
}

type RequestService struct {
	requestRepo models.RequestRepository
}

func NewRequestService(requestRepo models.RequestRepository) *RequestService {
	return &RequestService{
		requestRepo: requestRepo,
	}
}

func (s *RequestService) GetRequestById(ctx context.Context, id string) (*models.RequestData, error) {
	return s.requestRepo.FindByID(ctx, id)
}

func (s *RequestService) InsertManyRequests(ctx context.Context, requests []*models.RequestData) error {
	_, err := s.requestRepo.SaveAll(ctx, requests)
	return err
}

func (s *RequestService) SaveRequest(ctx context.Context, request *models.RequestData) (*models.RequestData, error) {
	return s.requestRepo.Save(ctx, request)
}

func (s *RequestService) GetPagedRequestsFilterdByUrl(ctx context.Context, url string, page, pageSize int) ([]*models.RequestData, error) {
	return s.requestRepo.FindAllByUrlLikePaged(ctx, url, page, pageSize)
}

func (s *RequestService) CountDocumentsForUrlFilter(ctx context.Context, url string) (int64, error) {
	return s.requestRepo.CountByUrlLike(ctx, url)
}
