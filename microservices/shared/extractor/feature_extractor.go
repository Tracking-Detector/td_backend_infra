package extractor

import (
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	"go.mongodb.org/mongo-driver/bson"
)

type ExtractorTypes int

const (
	DocumentId ExtractorTypes = iota
	DocumentLifecycle
	FrameId
	FrameType
	Initiator
	Method
	ParentFrameId
	RequestId
	TabId
	TimeStamp
	Type
	URL
	Success
	RequestHeaders
	Tracker
)

type DocumentIdExtractor func(string) ([]int, error)
type DocumentLifecycleExtractor func(string) ([]int, error)
type FrameIdExtractor func(int) ([]int, error)
type FrameTypeExtractor func(string) ([]int, error)
type InitiatorExtractor func(string) ([]int, error)
type MethodExtractor func(string) ([]int, error)
type ParentFrameIdExtractor func(int) ([]int, error)
type RequestIdExtractor func(string) ([]int, error)
type TabIdExtractor func(int) ([]int, error)
type TimeStampExtractor func(float32) ([]int, error)
type TypeExtractor func(string) ([]int, error)
type URLExtractor func(string) ([]int, error)
type SuccessExtractor func(string) ([]int, error)
type RequestHeadersExtractor func([]map[string]string) ([]int, error)
type TrackerExtractor func(bool) ([]int, error)

type Extractor struct {
	name                       string
	description                string
	dimensions                 []int
	query                      bson.M
	sequence                   []ExtractorTypes
	documentIdExtractor        DocumentIdExtractor
	documentLifecycleExtractor DocumentLifecycleExtractor
	frameIdExtractor           FrameIdExtractor
	frameTypeExtractor         FrameTypeExtractor
	initiatorExtractor         InitiatorExtractor
	methodExtractor            MethodExtractor
	parentFrameIdExtractor     ParentFrameIdExtractor
	requestIdExtractor         RequestIdExtractor
	tabIdExtractor             TabIdExtractor
	timeStampExtractor         TimeStampExtractor
	typeExtractor              TypeExtractor
	urlExtractor               URLExtractor
	successExtractor           SuccessExtractor
	requestHeadersExtractor    RequestHeadersExtractor
	trackerExtractor           TrackerExtractor
}

func NewExtractor(name string, description string, dimensions []int) *Extractor {
	return &Extractor{
		sequence:    make([]ExtractorTypes, 0),
		name:        name,
		dimensions:  dimensions,
		description: description,
		query:       bson.M{},
	}
}

func (e *Extractor) Query() bson.M {
	return e.query
}

func (e *Extractor) GetName() string {
	return e.name
}

func (e *Extractor) GetDimensions() []int {
	return e.dimensions
}

func (e *Extractor) GetDescription() string {
	return e.description
}

func (e *Extractor) GetFileName() string {
	return e.name + ".csv.gz"
}

func (e *Extractor) DocumentId(extractor DocumentIdExtractor) {
	e.documentIdExtractor = extractor
	e.sequence = append(e.sequence, DocumentId)
	e.query["documentId"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) DocumentLifecycle(extractor DocumentLifecycleExtractor) {
	e.documentLifecycleExtractor = extractor
	e.sequence = append(e.sequence, DocumentLifecycle)
	e.query["documentLifecycle"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) FrameId(extractor FrameIdExtractor) {
	e.frameIdExtractor = extractor
	e.sequence = append(e.sequence, FrameId)
	e.query["frameId"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) FrameType(extractor FrameTypeExtractor) {
	e.frameTypeExtractor = extractor
	e.sequence = append(e.sequence, FrameType)
	e.query["frameType"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) Initiator(extractor InitiatorExtractor) {
	e.initiatorExtractor = extractor
	e.sequence = append(e.sequence, Initiator)
	e.query["initiator"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) Method(extractor MethodExtractor) {
	e.methodExtractor = extractor
	e.sequence = append(e.sequence, Method)
	e.query["method"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) ParentFrameId(extractor ParentFrameIdExtractor) {
	e.parentFrameIdExtractor = extractor
	e.sequence = append(e.sequence, ParentFrameId)
	e.query["parentFrameId"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) RequestId(extractor RequestIdExtractor) {
	e.requestIdExtractor = extractor
	e.sequence = append(e.sequence, RequestId)
	e.query["requestId"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) TabId(extractor TabIdExtractor) {
	e.tabIdExtractor = extractor
	e.sequence = append(e.sequence, TabId)
	e.query["tabId"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) TimeStamp(extractor TimeStampExtractor) {
	e.timeStampExtractor = extractor
	e.sequence = append(e.sequence, TimeStamp)
	e.query["timeStamp"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) Type(extractor TypeExtractor) {
	e.typeExtractor = extractor
	e.sequence = append(e.sequence, Type)
	e.query["type"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) URL(extractor URLExtractor) {
	e.urlExtractor = extractor
	e.sequence = append(e.sequence, URL)
	e.query["url"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) Success(extractor SuccessExtractor) {
	e.successExtractor = extractor
	e.sequence = append(e.sequence, Success)
	e.query["success"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) RequestHeaders(extractor RequestHeadersExtractor) {
	e.requestHeadersExtractor = extractor
	e.sequence = append(e.sequence, RequestHeaders)
	e.query["requestHeaders"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) Tracker(extractor TrackerExtractor) {
	e.trackerExtractor = extractor
	e.sequence = append(e.sequence, Tracker)
	e.query["labels"] = bson.M{
		"$exists": true,
	}
}

func (e *Extractor) Encode(requestData models.ReducedRequestData) ([]int, error) {
	encoding := make([]int, 0)
	var err error
	var val []int
	for _, next := range e.sequence {
		switch next {
		case DocumentId:
			val, err = e.documentIdExtractor(requestData.DocumentId)
		case DocumentLifecycle:
			val, err = e.documentLifecycleExtractor(requestData.DocumentLifecycle)
		case FrameId:
			val, err = e.frameIdExtractor(requestData.FrameId)
		case FrameType:
			val, err = e.frameTypeExtractor(requestData.FrameType)
		case Initiator:
			val, err = e.initiatorExtractor(requestData.Initiator)
		case Method:
			val, err = e.methodExtractor(requestData.Method)
		case ParentFrameId:
			val, err = e.parentFrameIdExtractor(requestData.ParentFrameId)
		case RequestId:
			val, err = e.parentFrameIdExtractor(requestData.ParentFrameId)
		case TabId:
			val, err = e.tabIdExtractor(requestData.TabId)
		case TimeStamp:
			val, err = e.timeStampExtractor(requestData.TimeStamp)
		case Type:
			val, err = e.typeExtractor(requestData.Type)
		case URL:
			val, err = e.urlExtractor(requestData.URL)
		case Success:
			val, err = e.successExtractor(requestData.URL)
		case RequestHeaders:
			val, err = e.requestHeadersExtractor(requestData.RequestHeaders)
		case Tracker:
			val, err = e.trackerExtractor(requestData.Tracker)
		}
		if err != nil {
			return nil, err
		} else {
			encoding = append(encoding, val...)
		}
	}
	return encoding, nil
}
