package main

import "net/http"

func (app *application) NotFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sorry we don't have this page :("))
}
