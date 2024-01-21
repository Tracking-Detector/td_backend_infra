package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	ADMIN  Role = "admin"
	CLIENT Role = "client"
)

type UserData struct {
	Id    primitive.ObjectID `bson:"_id,omitempty"`
	Role  Role               `bson:"role"`
	Email string             `bson:"email"`
	Key   string             `bson:"key"`
}

type CreateUserData struct {
	Email string `bson:"email"`
}

type UserDataRepresentation struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Email string             `json:"email" bson:"email"`
	Role  Role               `json:"role" bson:"role"`
}
