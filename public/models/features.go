package models

type Feature struct {
	Title   string `json:"title" bson:"title"`
	Caption string `json:"caption" bson:"caption"`
}

type Features struct {
	Section string     `json:"section"  bson:"section"`
	Title   string     `json:"title"   bson:"title"`
	Caption string     `json:"caption" bson:"caption"`
	Items   []*Feature `json:"items"  bson:"items"`
}
