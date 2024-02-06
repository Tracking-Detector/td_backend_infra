package models

type Feature struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
}

type Features struct {
	Section string     `json:"section"`
	Title   string     `json:"title"`
	Caption string     `json:"caption"`
	Items   []*Feature `json:"items"`
}
