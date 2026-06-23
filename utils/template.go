package utils


import (
	"html/template"
)



var FuncMap = template.FuncMap{
	"formatLocation": FormattedDatesAndLocation,
}

var Tmpl = template.Must(
	template.New("").Funcs(FuncMap).ParseFiles(
		"templates/index.html",
		"templates/artist.html",
		"templates/home.html",
		"templates/search.html",
		"templates/error.html",
	),
)