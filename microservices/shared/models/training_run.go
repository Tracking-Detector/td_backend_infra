package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrainingRun struct {
	Id              primitive.ObjectID `json:"_id,omitempty"`
	Name            string             `json:"name"`
	DataSet         string             `json:"dataSet"`
	Time            string             `json:"time"`
	F1Train         float64            `json:"f1Train"`
	F1Test          float64            `json:"f1Test"`
	TrainingHistory bson.M             `json:"trainingHistory"`
	BatchSize       int                `json:"batchSize"`
	Epochs          int                `json:"epochs"`
}
