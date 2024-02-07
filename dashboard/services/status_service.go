package services

import (
	"sync"

	"github.com/Tracking-Detector/td_backend_infra/dashboard/config"
	"github.com/Tracking-Detector/td_backend_infra/dashboard/models"
)

type IStatusService interface {
	GetStatus() []*models.ServiceStatus
}

type StatusService struct {
	restService            IRestService
	serviceHealthEndpoints map[string]string
	cache                  []*models.ServiceStatus
}

// TODO add cron job to reload cache
func NewStatusService(restService IRestService) *StatusService {
	return &StatusService{
		restService:            restService,
		serviceHealthEndpoints: config.GetServicesHealth(),
	}
}

func (s *StatusService) ReloadCache() {
	var (
		wg              sync.WaitGroup
		mu              sync.Mutex
		serviceStatuses []*models.ServiceStatus
	)
	resultCh := make(chan *models.ServiceStatus, len(s.serviceHealthEndpoints))

	for serviceName, endpoint := range s.serviceHealthEndpoints {
		wg.Add(1)
		go func(name, ep string) {
			defer wg.Done()

			resp, err := s.restService.Get(ep)
			if err != nil {
				resultCh <- &models.ServiceStatus{
					ServiceName:    name,
					Status:         models.ERROR,
					StatusSubtitle: "Service unavailable",
					ResponseTime:   "N/A",
				}
				return
			}

			mu.Lock()
			defer mu.Unlock()
			if resp.Time() > 400 {
				resultCh <- &models.ServiceStatus{
					ServiceName:    name,
					Status:         models.WARNING,
					StatusSubtitle: "Response time slower than usual",
					ResponseTime:   resp.Time().String(),
				}
				return
			}
			resultCh <- &models.ServiceStatus{
				ServiceName:    name,
				Status:         models.HEALTHY,
				StatusSubtitle: "No issues detected",
				ResponseTime:   resp.Time().String(),
			}
		}(serviceName, endpoint)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for status := range resultCh {
		serviceStatuses = append(serviceStatuses, status)
	}
	s.cache = serviceStatuses
}

func (s *StatusService) GetStatus() []*models.ServiceStatus {
	if s.cache == nil || len(s.cache) == 0 {
		s.ReloadCache()
	}
	return s.cache
}
