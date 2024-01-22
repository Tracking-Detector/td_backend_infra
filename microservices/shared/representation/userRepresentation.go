package representation

import (
	"tds/shared/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDataRepresentation struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Email string             `json:"email" bson:"email"`
	Role  models.Role        `json:"role" bson:"role"`
}

func ConvertUserDataToUserDataRepresentation(user *models.UserData) *UserDataRepresentation {
	return &UserDataRepresentation{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}
}

func ConvertUserDatasToUserDataRepresentations(users []*models.UserData) []*UserDataRepresentation {
	userRepresentations := make([]*UserDataRepresentation, len(users))
	for i, u := range users {
		userRepresentations[i] = ConvertUserDataToUserDataRepresentation(u)
	}
	return userRepresentations
}
