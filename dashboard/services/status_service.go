package services

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Tracking-Detector/td_backend_infra/dashboard/config"
	"github.com/Tracking-Detector/td_backend_infra/dashboard/models"
)

type IStatusService interface {
	ReloadCache()
	GetStatus() []*models.ServiceStatus
}

type StatusService struct {
	restService            IRestService
	serviceHealthEndpoints map[string]string
	cache                  []*models.ServiceStatus
	lastUpdate             time.Time
}

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
	sort.Slice(serviceStatuses, func(i, j int) bool {
		return serviceStatuses[i].ServiceName < serviceStatuses[j].ServiceName
	})
	s.lastUpdate = time.Now()
	s.cache = serviceStatuses
}

func (s *StatusService) GetStatus() []*models.ServiceStatus {
	if s.cache == nil || len(s.cache) == 0 || time.Since(s.lastUpdate) > 3*time.Hour {
		s.ReloadCache()
	}
	for _, status := range s.cache {
		status.LastUpdate = s.GetUpdateSinceString()
	}
	return s.cache
}

func (s *StatusService) GetUpdateSinceString() string {
	if s.lastUpdate.IsZero() {
		return "N/A"
	}
	duration := time.Since(s.lastUpdate)
	minutes := int(duration.Minutes())

	switch {
	case minutes < 1:
		return "Just now"
	case minutes < 60:
		return fmt.Sprintf("%d minutes ago", minutes)
	case minutes < 120:
		return "1 hour ago"
	case minutes < 180:
		return "2 hours ago"
	default:
		hours := minutes / 60
		minutes %= 60
		return fmt.Sprintf("%d hours and %d minutes ago", hours, minutes)
	}
}
