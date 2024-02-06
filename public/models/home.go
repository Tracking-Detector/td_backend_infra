package models

type Home struct {
	Title    string
	Hero     *Hero
	Features *Features
	Products *Products
}

func NewHome() *Home {
	return &Home{
		Title: "Tracker Detector | Home",
	}
}
