package service

import (
	"context"
	"errors"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/configs"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	log "github.com/sirupsen/logrus"
)

type IUserService interface {
	GetAllUsers(ctx context.Context) ([]*models.UserData, error)
	DeleteUserByID(ctx context.Context, id string) error
	InitAdmin(ctx context.Context) error
	CreateApiUser(ctx context.Context, email string) (string, error)
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
	err := s.userRepository.InTransaction(ctx, func(ctx context.Context) error {
		user, err := s.userRepository.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if user.Role == models.ADMIN {
			return errors.New("cannot delete admin users")
		}
		return s.userRepository.DeleteByID(ctx, id)
	})
	return err
}

func (s *UserService) InitAdmin(ctx context.Context) error {
	adminApiKey := configs.EnvAdminApiKey()
	hash, _ := s.encryptionService.HashPassword(adminApiKey)
	admin := &models.UserData{
		Role:  models.ADMIN,
		Email: configs.EnvAdminEmail(),
		Key:   string(hash),
	}
	_, err := s.userRepository.Save(ctx, admin)
	return err
}

func (s *UserService) CreateApiUser(ctx context.Context, email string) (string, error) {
	hash, pw, err := s.encryptionService.GenerateApiKey()
	if err != nil {
		log.WithFields(log.Fields{
			"service": "UserService",
			"error":   err.Error(),
		}).Fatal("Could not generate Hash and ApiKey")
		return "", err
	}
	err = s.userRepository.InTransaction(ctx, func(ctx context.Context) error {
		user, _ := s.userRepository.FindByEmail(ctx, email)
		if user != nil {
			return errors.New("user with email already registered")
		}

		newUser := &models.UserData{
			Role:  models.CLIENT,
			Email: email,
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
