package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	Id                primitive.ObjectID  `json:"_id" bson:"_id,omitempty"`
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
	Labels            []RequestDataLabel  `json:"labels" bson:"labels" validate:"required"`
}
