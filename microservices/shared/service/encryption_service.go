package service

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type IEncryptionService interface {
	GenerateApiKey() (string, string, error)
	HashPassword(password string) (string, error)
}

type EncryptionService struct {
}

func NewEncryptionService() *EncryptionService {
	return &EncryptionService{}
}

func (s *EncryptionService) GenerateApiKey() (string, string, error) {
	key := uuid.New().String()
	hash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	return string(hash), key, err
}

func (s *EncryptionService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"service": "EncryptionService",
			"error":   err.Error(),
		}).Fatal("Could not generate Hash for Password")
	}
	return string(hash), nil
}
