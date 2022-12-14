package main

import (
	"embed"
	"encoding/json"
	"net/http"
)

//go:embed markers.json
var markers embed.FS

func (app *application) NotFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sorry we don't have this page :("))
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// Get the markers info

	b, err := markers.ReadFile("markers.json")
	if err != nil {
		app.serverError(w, err)
		return
	}
	var markers Markers
	err = json.Unmarshal(b, &markers)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.go.tmpl", &templateData{Markers: markers.Markers})
}
