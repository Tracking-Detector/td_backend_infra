package responses

import "tds/shared/models"

type ModelDataResponse struct {
	Status int                `json:"status"`
	Data   []models.ModelData `json:"data"`
}
