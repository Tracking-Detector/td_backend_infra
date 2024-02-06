package models

type Product struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Logo        string   `json:"logo"`
	BulletPoint []string `json:"points"`
}

type Products struct {
	Section  string     `json:"section"`
	Title    string     `json:"title"`
	Caption  string     `json:"caption"`
	Products []*Product `json:"products"`
}
