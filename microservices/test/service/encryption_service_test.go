package service_test

import (
	"tds/shared/service"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

func TestEncryptionService(t *testing.T) {
	suite.Run(t, &EncryptionServiceTest{})
}

type EncryptionServiceTest struct {
	suite.Suite
	encrytpionService *service.EncryptionService
}

func (suite *EncryptionServiceTest) SetupTest() {
	suite.encrytpionService = service.NewEncryptionService()
}

func (suite *EncryptionServiceTest) TestShouldGenerateApiKey() {
	// given

	// when
	hash, key, err := suite.encrytpionService.GenerateApiKey()
	// then
	suite.Nil(err)
	suite.NotEmpty(hash)
	suite.NotEmpty(key)
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(key))
	suite.Nil(err)
}

func (suite *EncryptionServiceTest) TestShouldHashPassword() {
	// given
	pw := "somePw"
	// when
	hash, err := suite.encrytpionService.HashPassword(pw)
	// then
	suite.Nil(err)
	suite.NotEmpty(hash)
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	suite.Nil(err)
}

func (suite *EncryptionServiceTest) TestShouldCompareHashAndPassword() {
	// given
	pw := "somePw"
	hash, err := suite.encrytpionService.HashPassword(pw)
	suite.Nil(err)
	// when
	valid := suite.encrytpionService.CompareHashAndPassword(hash, pw)
	// then
	suite.True(valid)

}
