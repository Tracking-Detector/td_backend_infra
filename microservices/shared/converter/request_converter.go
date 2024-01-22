package converter

import (
	"tds/shared/models"
)

type ReduceType string
type ReduceFunction func([]models.RequestDataLabel) bool

const (
	OR           ReduceType = "or"
	AND          ReduceType = "and"
	EASY_PRIVACY ReduceType = "EasyPrivacy"
	EASY_LIST    ReduceType = "EasyList"
	HUMAN        ReduceType = "Human"
)

func orFunction(labels []models.RequestDataLabel) bool {
	isTracking := false
	for _, label := range labels {
		isTracking = isTracking || label.IsLabeled
	}
	if isTracking {
		return true
	}
	return false
}

func andFunction(labels []models.RequestDataLabel) bool {
	isTracking := true
	for _, label := range labels {
		isTracking = isTracking && label.IsLabeled
	}
	if isTracking {
		return true
	}
	return false
}

func easyPrivacyFunction(labels []models.RequestDataLabel) bool {
	for _, label := range labels {
		if label.Blocklist == "EasyPrivacy" {
			if label.IsLabeled {
				return true
			} else {
				return false
			}
		}

	}
	return false
}

func easyListFunction(labels []models.RequestDataLabel) bool {
	for _, label := range labels {
		if label.Blocklist == "EasyList" {
			if label.IsLabeled {
				return true
			} else {
				return false
			}
		}

	}
	return false
}

func humanFunction(labels []models.RequestDataLabel) bool {
	for _, label := range labels {
		if label.Blocklist == "Human" {
			if label.IsLabeled {
				return true
			} else {
				return false
			}
		}

	}
	return false
}

func ConvertRequestModel(request *models.RequestData, reducer ReduceType) *models.ReducedRequestData {
	if request == nil {
		return nil
	}
	if len(request.Labels) == 0 {
		return nil
	}

	var reducerFn ReduceFunction
	switch reducer {
	case OR:
		reducerFn = orFunction
	case AND:
		reducerFn = andFunction
	case EASY_PRIVACY:
		reducerFn = easyPrivacyFunction
	case EASY_LIST:
		reducerFn = easyListFunction
	case HUMAN:
		reducerFn = humanFunction
	}

	reducedModel := &models.ReducedRequestData{}
	reducedModel.DocumentId = request.DocumentId
	reducedModel.DocumentLifecycle = request.DocumentLifecycle
	reducedModel.FrameId = request.FrameId
	reducedModel.FrameType = request.FrameType
	reducedModel.Initiator = request.Initiator
	reducedModel.Method = request.Method
	reducedModel.ParentFrameId = request.ParentFrameId
	reducedModel.RequestHeaders = request.RequestHeaders
	reducedModel.RequestId = request.RequestId
	reducedModel.Response = request.Response
	reducedModel.Success = request.Success
	reducedModel.TabId = request.TabId
	reducedModel.TimeStamp = request.TimeStamp
	reducedModel.Type = request.Type
	reducedModel.URL = request.URL
	reducedModel.Tracker = reducerFn(request.Labels)
	return reducedModel
}
