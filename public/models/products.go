package models

type Product struct {
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	Logo        string   `json:"logo" bson:"logo"`
	BulletPoint []string `json:"points" bson:"points"`
}

type Products struct {
	ID       string     `json:"id,omitempty" bson:"_id,omitempty"`
	Section  string     `json:"section" bson:"section"`
	Title    string     `json:"title" bson:"title"`
	Caption  string     `json:"caption" bson:"caption"`
	Products []*Product `json:"products" bson:"products"`
}
