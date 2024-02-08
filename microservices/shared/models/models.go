package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type BaseModel interface {
	GetID() string
	SetID(id string)
}

type BaseModelName interface {
	BaseModel
	GetName() string
}

// Exporter
type RunType string

const (
	IN_SERVICE RunType = "in-service"
	JS         RunType = "js"
)

type Exporter struct {
	ID                   string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name                 string  `json:"name" bson:"name"`
	Description          string  `json:"description" bson:"description"`
	Dimensions           []int   `json:"dimensions" bson:"dimensions"`
	Type                 RunType `json:"type" bson:"type"`
	ExportScriptLocation *string `json:"location" bson:"location"`
}

func (e *Exporter) GetID() string {
	return e.ID
}

func (e *Exporter) SetID(id string) {
	e.ID = id
}

func (e *Exporter) GetName() string {
	return e.Name
}

// Model
type Model struct {
	ID               string `json:"id" bson:"_id"`
	Name             string `json:"name" bson:"name"`
	Description      string `json:"description" bson:"description"`
	Dims             []int  `json:"dims" bson:"dims"`
	TensorflowLayers bson.D `json:"tfLayers" bson:"tfLayers"`
}

func (e *Model) GetID() string {
	return e.ID
}

func (e *Model) SetID(id string) {
	e.ID = id
}

func (e *Model) GetName() string {
	return e.Name
}

// TrainingRun
type TrainingRun struct {
	ID              string  `json:"id,omitempty" bson:"_id,omitempty"`
	ModelId         string  `json:"modelId" bson:"modelId"`
	Name            string  `json:"name" bson:"name"`
	DataSet         string  `json:"dataSet" bson:"dataSet"`
	Time            string  `json:"time" bson:"time"`
	F1Train         float64 `json:"f1Train" bson:"f1Train"`
	F1Test          float64 `json:"f1Test" bson:"f1Test"`
	TrainingHistory bson.M  `json:"trainingHistory" bson:"trainingHistory"`
	BatchSize       int     `json:"batchSize" bson:"batchSize"`
	Epochs          int     `json:"epochs" bson:"epochs"`
}

func (e *TrainingRun) GetID() string {
	return e.ID
}

func (e *TrainingRun) SetID(id string) {
	e.ID = id
}

func (e *TrainingRun) GetName() string {
	return e.Name
}

// User
type Role string

const (
	ADMIN  Role = "admin"
	CLIENT Role = "client"
)

type UserData struct {
	ID    string `json:"id,omitempty" bson:"_id,omitempty"`
	Role  Role   `json:"role" bson:"role"`
	Email string `json:"email" bson:"email"`
	Key   string `json:"key" bson:"key"`
}

func (e *UserData) GetID() string {
	return e.ID
}

func (e *UserData) SetID(id string) {
	e.ID = id
}

// Request
type RequestDataLabel struct {
	IsLabeled bool   `json:"isLabeled"  validate:"required"`
	Blocklist string `json:"blocklist"  validate:"required"`
}

type RequestDataResponse struct {
	DocumentId        string              `json:"documentId"`
	DocumentLifecycle string              `json:"documentLifecycle"`
	FrameId           int                 `json:"frameId"`
	FrameType         string              `json:"frameType"`
	FromCache         bool                `json:"fromCache"`
	Initiator         string              `json:"initiator"`
	Ip                string              `json:"ip"`
	Method            string              `json:"method"`
	ParentFrameId     int                 `json:"parentFrameId"`
	RequestId         string              `json:"requestId"`
	RequestHeaders    []map[string]string `json:"responseHeaders"`
	StatusCode        int                 `json:"statusCode"`
	StatusLine        string              `json:"statusLine"`
	TabId             int                 `json:"tabId"`
	TimeStamp         float32             `json:"timeStamp"`
	Type              string              `json:"type"`
	URL               string              `json:"url"`
}

type RequestData struct {
	ID                string              `json:"id" bson:"_id,omitempty"`
	DocumentId        string              `json:"documentId" bson:"documentId"`
	DocumentLifecycle string              `json:"documentLifecycle" bson:"documentLifecycle"`
	FrameId           int                 `json:"frameId" bson:"frameId"`
	FrameType         string              `json:"frameType" bson:"frameType"`
	Initiator         string              `json:"initiator" bson:"initiator"`
	Method            string              `json:"method" bson:"method"`
	ParentFrameId     int                 `json:"parentFrameId" bson:"parentFrameId"`
	RequestId         string              `json:"requestId" bson:"requestId"`
	TabId             int                 `json:"tabId" bson:"tabId"`
	TimeStamp         float32             `json:"timeStamp" bson:"timeStamp"`
	Type              string              `json:"type" bson:"type"`
	URL               string              `json:"url" bson:"url" validate:"required"`
	RequestHeaders    []map[string]string `json:"requestHeaders" bson:"requestHeaders"`
	Response          RequestDataResponse `json:"response" bson:"response"`
	Success           bool                `json:"success" bson:"success"`
	Dataset           string              `json:"dataset" bson:"dataset"`
	Labels            []RequestDataLabel  `json:"labels" bson:"labels" validate:"required"`
}

func (e *RequestData) GetID() string {
	return e.ID
}

func (e *RequestData) SetID(id string) {
	e.ID = id
}

type ReducedRequestData struct {
	DocumentId        string              `json:"documentId" bson:"documentId"`
	DocumentLifecycle string              `json:"documentLifecycle" bson:"documentLifecycle"`
	FrameId           int                 `json:"frameId" bson:"frameId"`
	FrameType         string              `json:"frameType" bson:"frameType"`
	Initiator         string              `json:"initiator" bson:"initiator"`
	Method            string              `json:"method" bson:"method"`
	ParentFrameId     int                 `json:"parentFrameId" bson:"parentFrameId"`
	RequestId         string              `json:"requestId" bson:"requestId"`
	TabId             int                 `json:"tabId" bson:"tabId"`
	TimeStamp         float32             `json:"timeStamp" bson:"timeStamp"`
	Type              string              `json:"type" bson:"type"`
	URL               string              `json:"url" bson:"url" validate:"required"`
	RequestHeaders    []map[string]string `json:"requestHeaders" bson:"requestHeaders"`
	Response          RequestDataResponse `json:"response" bson:"response"`
	Success           bool                `json:"success" bson:"success"`
	Tracker           bool                `json:"tracker" bson:"tracker"`
}

// Export Metrics
type ExportMetrics struct {
	Tracker    int    `json:"tracker" bson:"tracker"`
	NonTracker int    `json:"nonTracker" bson:"nonTracker"`
	Total      int    `json:"total" bson:"total"`
	Error      string `json:"error" bson:"error"`
}

type ExportRun struct {
	ID         string         `json:"id" bson:"_id,omitempty"`
	ExporterId string         `json:"exporterId" bson:"exporterId"`
	Name       string         `json:"name" bson:"name"`
	Reducer    string         `json:"reducer" bson:"reducer"`
	Dataset    string         `json:"dataset" bson:"dataset"`
	Metrics    *ExportMetrics `json:"metrics" bson:"metrics"`
	Start      time.Time      `json:"start" bson:"start"`
	End        time.Time      `json:"end" bson:"end"`
}

func (e *ExportRun) GetID() string {
	return e.ID
}

func (e *ExportRun) SetID(id string) {
	e.ID = id
}

func (e *ExportRun) GetName() string {
	return e.Name
}

// Dataset

type ReducerMetric struct {
	Reducer    string `json:"reducer" bson:"reducer"`
	Total      int    `json:"total" bson:"total"`
	Tracker    int    `json:"tracker" bson:"tracker"`
	NonTracker int    `json:"nonTracker" bson:"nonTracker"`
}
type DataSetMetrics struct {
	Total         int              `json:"total" bson:"total"`
	ReducerMetric []*ReducerMetric `json:"reducerMetric" bson:"reducerMetric"`
}

type Dataset struct {
	ID          string          `json:"id" bson:"_id,omitempty"`
	Name        string          `json:"name" bson:"name"`
	Label       string          `json:"label,omitempty" bson:"label,omitempty"`
	Description string          `json:"description" bson:"description"`
	Metrics     *DataSetMetrics `json:"metrics,omitempty" bson:"metrics,omitempty"`
}

func (e *Dataset) GetID() string {
	return e.ID
}

func (e *Dataset) SetID(id string) {
	e.ID = id
}

func (e *Dataset) GetName() string {
	return e.Name
}
