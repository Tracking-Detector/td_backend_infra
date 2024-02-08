package models

type InstallationGuide struct {
	ID      string `json:"id,omitempty" bson:"_id,omitempty"`
	Title   string `json:"title" bson:"title"`
	Caption string `json:"caption" bson:"caption"`
}
