package models

import "github.com/a-h/templ"

type Link struct {
	Text string
	URL  templ.SafeURL
}

type Navbar struct {
	Links []Link
}
