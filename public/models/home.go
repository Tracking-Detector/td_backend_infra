package models

type Home struct {
	Title             string
	Hero              *Hero
	Features          *Features
	Products          *Products
	InstallationGuide *InstallationGuide
}

func NewHome() *Home {
	return &Home{
		Title: "Tracker Detector | Home",
	}
}
