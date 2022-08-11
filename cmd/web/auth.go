package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost/4001/callback",
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		// What permissions we want to ask from the user.
		Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	// random string that identifies the user's request.
	randomState = "random"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body> <a href="/user/login"> Google Login </a> </body> </html>`
	fmt.Fprint(w, html)
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	app.infoLog.Println("url", url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback gets a respond from google and try to get the user info.
func (app *application) Callback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != randomState {
		app.errorLog.Print("State is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		app.errorLog.Printf("Could not get token: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	app.infoLog.Println("access token: ", token.AccessToken)
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		app.errorLog.Printf("Could not create get request: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.errorLog.Printf("Could not parse response: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}
	fmt.Fprintf(w, string(body))
}
