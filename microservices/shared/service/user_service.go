package service

import (
	"context"
	"errors"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/payload"
	log "github.com/sirupsen/logrus"
)

type IUserService interface {
	GetAllUsers(ctx context.Context) ([]*models.UserData, error)
	DeleteUserByID(ctx context.Context, id string) error
	CreateApiUser(ctx context.Context, email payload.CreateUserData) (string, error)
}

type UserService struct {
	userRepository    models.UserRepository
	encryptionService IEncryptionService
}

func NewUserService(userRepository models.UserRepository, encryptionService IEncryptionService) *UserService {
	return &UserService{
		userRepository:    userRepository,
		encryptionService: encryptionService,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.UserData, error) {
	return s.userRepository.FindAll(ctx)
}

func (s *UserService) DeleteUserByID(ctx context.Context, id string) error {
	return s.userRepository.DeleteByID(ctx, id)
}

func (s *UserService) CreateApiUser(ctx context.Context, userPayload payload.CreateUserData) (string, error) {
	hash, pw, err := s.encryptionService.GenerateApiKey()
	if err != nil {
		log.WithFields(log.Fields{
			"service": "UserService",
			"error":   err.Error(),
		}).Fatal("Could not generate Hash and ApiKey")
		return "", err
	}
	err = s.userRepository.InTransaction(ctx, func(ctx context.Context) error {
		user, _ := s.userRepository.FindByEmail(ctx, userPayload.Email)
		if user != nil {
			return errors.New("user with email already registered")
		}

		newUser := &models.UserData{
			Role:  userPayload.Role,
			Email: userPayload.Email,
			Key:   hash,
		}
		_, err := s.userRepository.Save(ctx, newUser)
		return err
	})
	if err != nil {
		return "", err
	}
	return pw, nil
}
