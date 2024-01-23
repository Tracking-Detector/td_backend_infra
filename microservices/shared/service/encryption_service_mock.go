package service

import "github.com/stretchr/testify/mock"

type EncryptionServiceMock struct {
	mock.Mock
}

func NewEncryptionServiceMock() *EncryptionServiceMock {
	return &EncryptionServiceMock{}
}

func (m *EncryptionServiceMock) GenerateApiKey() (string, string, error) {
	args := m.Called()
	return args.String(0), args.String(1), args.Error(2)
}

func (m *EncryptionServiceMock) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}
