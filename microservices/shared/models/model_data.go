package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ModelData struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Dims        []int              `json:"dims" bson:"dims"`
}
