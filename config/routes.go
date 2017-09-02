package config

import (
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/controllers"
	"github.com/omar-ozgur/flock-api/middleware"
	"net/http"
)

// router holds application router information
type router struct {
	muxRouter  *mux.Router
	middleware middleware.Middleware
}

// init initializes the router
func (r *router) init() {

	// Create routes
	r.muxRouter = mux.NewRouter()
	r.createRoutes()

	// Initialize middleware
	r.middleware = middleware.Middleware{}
	r.middleware.Init(r.muxRouter)
}

// handle registers a path and handler with the router
func (r *router) handle(path string, handler http.Handler) *mux.Route {
	return r.muxRouter.Handle(path, handler)
}

// authorizationHandler runs authorization middleware for the specified handler
func (r *router) authorizationHandler(handler http.Handler) http.Handler {
	return middleware.JWTMiddleware.Handler(handler)
}

// createRoutes creates application routes
func (r *router) createRoutes() {
	r.createUserRoutes()
	r.createEventRoutes()
}

// createUserRoutes creates user routes
func (r *router) createUserRoutes() {

	// Signup a new user
	r.handle(
		"/signup",
		controllers.UsersCreate,
	).Methods("POST")

	// Login an existing user
	r.handle(
		"/login",
		controllers.UsersLogin,
	).Methods("POST")

	// View the current user
	r.handle(
		"/profile",
		r.authorizationHandler(controllers.UsersProfile),
	).Methods("Get")

	// View all users
	r.handle(
		"/users",
		controllers.UsersIndex,
	).Methods("GET")

	// Search for a user
	r.handle(
		"/users/search",
		controllers.UsersSearch,
	).Methods("POST")

	// View a specific user
	r.handle(
		"/users/{id}",
		controllers.UsersShow,
	).Methods("GET")

	// Update a specific user
	r.handle(
		"/users/{id}",
		r.authorizationHandler(controllers.UsersUpdate),
	).Methods("PUT")

	// Delete a specific user
	r.handle(
		"/users/{id}",
		r.authorizationHandler(controllers.UsersDelete),
	).Methods("DELETE")

	// See the current user's attendance
	r.handle(
		"/users/{id}/attendance",
		r.authorizationHandler(controllers.UsersAttendance),
	).Methods("GET")

	// Login with facebook
	r.handle(
		"/loginWithFacebook",
		controllers.LoginWithFacebook,
	).Methods("POST")
}

// createEventRoutes creates event routes
func (r *router) createEventRoutes() {

	// View all events
	r.handle(
		"/events",
		controllers.EventsIndex,
	).Methods("GET")

	// Create an event
	r.handle(
		"/events",
		r.authorizationHandler(controllers.EventsCreate),
	).Methods("POST")

	// Search for an event
	r.handle(
		"/events/search",
		controllers.EventsSearch,
	).Methods("POST")

	// View a specific event
	r.handle(
		"/events/{id}",
		controllers.EventsShow,
	).Methods("GET")

	// Update a specific event
	r.handle(
		"/events/{id}",
		r.authorizationHandler(controllers.EventsUpdate),
	).Methods("PUT")

	// Delete a specific event
	r.handle(
		"/events/{id}",
		r.authorizationHandler(controllers.EventsDelete),
	).Methods("DELETE")

	// View the attendees of a specific event
	r.handle(
		"/events/{id}/attendees",
		controllers.EventsAttendees,
	).Methods("GET")

	// Have the current user attend a specific event
	r.handle(
		"/events/{id}/attend",
		r.authorizationHandler(controllers.EventsAttend),
	).Methods("POST")

	// Remove the current user for a specific event's attendees
	r.handle(
		"/events/{id}/attendance",
		r.authorizationHandler(controllers.EventsDeleteAttendance),
	).Methods("DELETE")
}
