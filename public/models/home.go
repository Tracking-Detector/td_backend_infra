package models

type Home struct {
	Title             string
	Navbar            *Navbar
	Hero              *Hero
	Features          *Features
	Products          *Products
	InstallationGuide *InstallationGuide
	Contact           *Contact
}

func NewHome() *Home {
	return &Home{
		Title: "Tracker Detector | Home",
	}
}
