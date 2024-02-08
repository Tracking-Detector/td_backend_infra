package payload

import "github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"

type CreateUserData struct {
	Email string      `bson:"email"`
	Role  models.Role `bson:"role"`
}
