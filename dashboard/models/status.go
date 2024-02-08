package models

const (
	HEALTHY string = "Healthy"
	WARNING string = "Warning"
	ERROR   string = "Error"
)

type ServiceStatus struct {
	Status         string `json:"status"`
	StatusSubtitle string `json:"status_subtitle"`
	ServiceName    string `json:"service_name"`
	ResponseTime   string `json:"response_time"`
	LastUpdate     string `json:"last_update"`
}
