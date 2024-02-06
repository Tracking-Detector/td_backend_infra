package models

type Hero struct {
	Title    string `json:"title" bson:"title"`
	SubTitle string `json:"subTitle" bson:"subTitle"`
	Caption  string `json:"caption" bson:"caption"`
	Logo     string `json:"logo" bson:"logo"`
}
