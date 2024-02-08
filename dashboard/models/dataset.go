package models

type ReducerMetric struct {
	Reducer    string `json:"reducer" bson:"reducer"`
	Total      int    `json:"total" bson:"total"`
	Tracker    int    `json:"tracker" bson:"tracker"`
	NonTracker int    `json:"nonTracker" bson:"nonTracker"`
}
type DataSetMetrics struct {
	Total         int              `json:"total" bson:"total"`
	ReducerMetric []*ReducerMetric `json:"reducerMetric" bson:"reducerMetric"`
}

type Dataset struct {
	ID          string          `json:"id" bson:"_id,omitempty"`
	Name        string          `json:"name" bson:"name"`
	Label       string          `json:"label" bson:"label"`
	Description string          `json:"description" bson:"description"`
	Metrics     *DataSetMetrics `json:"metrics,omitempty" bson:"metrics,omitempty"`
}

type CreateDatasetPayload struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
	Label       string `json:"label" form:"label"`
}
