package converter_test

import (
	"testing"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/converter"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	"github.com/stretchr/testify/suite"
)

func TestRequestConverter(t *testing.T) {
	suite.Run(t, &TestRequestConverterTest{})
}

type TestRequestConverterTest struct {
	suite.Suite
}

func (suite *TestRequestConverterTest) TestConvertRequestModel_NilOnNilRequest() {
	// given

	// when
	request := converter.ConvertRequestModel(nil, converter.OR)
	// then
	suite.Nil(request)
}

func (suite *TestRequestConverterTest) TestConvertRequestModel_SuccessHumanReducer() {
	// given
	request := &models.RequestData{
		ID: "id",
		Labels: []models.RequestDataLabel{
			{
				IsLabeled: true,
				Blocklist: "Human",
			},
			{
				IsLabeled: false,
				Blocklist: "EasyPrivacy",
			},
		},
	}
	// when
	converted := converter.ConvertRequestModel(request, converter.HUMAN)

	// then
	suite.True(converted.Tracker)
}
