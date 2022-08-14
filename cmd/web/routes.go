package main

import (
	"net/http"

	"github.com/bmizerany/pat"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Create a new middleware chain containing the middleware specific to
	// our *dynamic application routes*. For now, this chain will only contain
	// the session middleware but we'll add more to it later.
	// Use the nosurf middleware on all our 'dynamic' routes.
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	// Swap the route declarations to use the application struct's methods as the
	// handler functions.
	// Update these routes to use the new dynamic middleware chain followed
	// by the appropriate handler function.
	mux := pat.New()
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.Login))
	mux.Get("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.Logout))
	mux.Get("/display", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.Display))

	mux.Get("/", dynamicMiddleware.ThenFunc(app.Home))

	mux.Get("/callback", dynamicMiddleware.ThenFunc(app.Callback))

	mux.NotFound = http.HandlerFunc(app.NotFound)

	// Return the 'standard' middleware chain followed by the servemux.
	return standardMiddleware.Then(mux)
}
