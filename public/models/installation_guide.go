package models

type InstallationGuide struct {
	Title   string `json:"title" bson:"title"`
	Caption string `json:"caption" bson:"caption"`
}
