package models

import (
	"encoding/json"
	"errors"
)

type JobPayload struct {
	FunctionName string   `json:"functionName"`
	Args         []string `json:"args"`
}

func NewJob(functionName string, args []string) *JobPayload {
	return &JobPayload{
		FunctionName: functionName,
		Args:         args,
	}
}

func (j *JobPayload) Serialize() (string, error) {
	data, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func DeserializeJob(data string) (*JobPayload, error) {
	var job JobPayload
	err := json.Unmarshal([]byte(data), &job)
	if err != nil {
		return nil, err
	}
	if job.FunctionName == "" {
		return nil, errors.New("invalid job data")
	}
	return &job, nil
}
