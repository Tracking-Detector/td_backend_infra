package service

import (
	"context"
	"errors"
	"strings"
	"tds/shared/models"
)

type IAuthService interface {
	ValidateBearerToken(token string) (bool, error)
}

type AuthService struct {
	userRepo          models.UserRepository
	encryptionService IEncryptionService
}

func NewAuthService(userRepo models.UserRepository, encryptionService IEncryptionService) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		encryptionService: encryptionService,
	}
}

func (s *AuthService) ValidateBearerToken(ctx context.Context, token string, isAdmin bool) (bool, error) {
	split := strings.Split(token, " ")
	if len(split) != 2 || split[0] != "Bearer" {
		return false, errors.New("Wrong formatted Bearer Token in Auth header.")
	}

	users, err := s.userRepo.FindAll(ctx)

	if err != nil {
		return false, err
	}
	for _, u := range users {
		if valid := s.encryptionService.CompareHashAndPassword(u.Key, split[1]); valid {
			if isAdmin {
				if u.Role == models.ADMIN {
					return true, nil
				}
			} else {
				return true, nil
			}
		}
	}
	return false, nil

}
