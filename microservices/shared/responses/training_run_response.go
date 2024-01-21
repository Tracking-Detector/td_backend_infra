package responses

import (
	"tds/shared/models"
)

type TrainingRunResponse struct {
	Status int                  `json:"status"`
	Data   []models.TrainingRun `json:"data"`
}
