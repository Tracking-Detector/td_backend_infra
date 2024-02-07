package models

type Hero struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Title    string `json:"title" bson:"title"`
	SubTitle string `json:"subTitle" bson:"subTitle"`
	Caption  string `json:"caption" bson:"caption"`
	Logo     string `json:"logo" bson:"logo"`
}
