package main

import (
	"net/http"
)

func (app *application) NotFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sorry we don't have this page :("))
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// Use the new render helper.
	app.render(w, r, "home.page.go.tmpl", &templateData{})
}
