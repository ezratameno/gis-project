package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// UserInfo is the user information we get from google
type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:4001/callback",
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		// What permissions we want to ask from the user.
		Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	// random string that identifies the user's request.
	randomState = "random"
)

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "user")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

// Callback gets a respond from google and try to get the user info.
func (app *application) Callback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != randomState {
		app.errorLog.Println("State is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		app.errorLog.Printf("code exchange failed: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		app.errorLog.Printf("failed getting user info: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		app.errorLog.Printf("failed reading response body: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	app.infoLog.Println(string(contents))

	var user UserInfo
	err = json.Unmarshal(contents, &user)
	if err != nil {
		app.errorLog.Printf("failed unmarshalling the content : %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// Add the user to the session, so that they are now 'logged
	// in'.
	app.session.Put(r, "user", user)
	app.session.Put(r, "flash", "You've been logged in successfully!")
	http.Redirect(w, r, "/", http.StatusOK)

}
