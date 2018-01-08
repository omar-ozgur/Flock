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
	middleware *middleware.Middleware
}

// NewRouter initializes a new router
func NewRouter() *router {

	router := router{}

	// Create routes
	router.muxRouter = mux.NewRouter()
	router.createRoutes()

	// Initialize middleware
	router.middleware = middleware.NewMiddleware(router.muxRouter)

	return &router
}

// handle registers a path and handler with the router
func (router *router) handle(path string, handler http.Handler) *mux.Route {
	return router.muxRouter.Handle(path, handler)
}

// authorizationHandler runs authorization middleware for the specified handler
func (router *router) authorizationHandler(handler http.Handler) http.Handler {
	return middleware.JWTMiddleware.Handler(handler)
}

// createRoutes creates application routes
func (router *router) createRoutes() {
	router.createUserRoutes()
	router.createEventRoutes()
}

// createUserRoutes creates user routes
func (router *router) createUserRoutes() {

	// Signup a new user
	router.handle(
		"/signup",
		controllers.UsersCreate,
	).Methods("POST")

	// Login an existing user
	router.handle(
		"/login",
		controllers.UsersLogin,
	).Methods("POST")

	// View the current user
	router.handle(
		"/profile",
		router.authorizationHandler(controllers.UsersProfile),
	).Methods("Get")

	// View all users
	router.handle(
		"/users",
		controllers.UsersIndex,
	).Methods("GET")

	// Search for a user
	router.handle(
		"/users/search",
		controllers.UsersSearch,
	).Methods("POST")

	// View a specific user
	router.handle(
		"/users/{id}",
		controllers.UsersShow,
	).Methods("GET")

	// Update a specific user
	router.handle(
		"/users/{id}",
		router.authorizationHandler(controllers.UsersUpdate),
	).Methods("PUT")

	// Delete a specific user
	router.handle(
		"/users/{id}",
		router.authorizationHandler(controllers.UsersDelete),
	).Methods("DELETE")

	// See the current user's attendance
	router.handle(
		"/users/{id}/attendance",
		router.authorizationHandler(controllers.UsersAttendance),
	).Methods("GET")

	// Login with facebook
	router.handle(
		"/loginWithFacebook",
		controllers.LoginWithFacebook,
	).Methods("POST")
}

// createEventRoutes creates event routes
func (router *router) createEventRoutes() {

	// View all events
	router.handle(
		"/events",
		controllers.EventsIndex,
	).Methods("GET")

	// Create an event
	router.handle(
		"/events",
		router.authorizationHandler(controllers.EventsCreate),
	).Methods("POST")

	// Search for an event
	router.handle(
		"/events/search",
		controllers.EventsSearch,
	).Methods("POST")

	// View a specific event
	router.handle(
		"/events/{id}",
		controllers.EventsShow,
	).Methods("GET")

	// Update a specific event
	router.handle(
		"/events/{id}",
		router.authorizationHandler(controllers.EventsUpdate),
	).Methods("PUT")

	// Delete a specific event
	router.handle(
		"/events/{id}",
		router.authorizationHandler(controllers.EventsDelete),
	).Methods("DELETE")

	// View the attendees of a specific event
	router.handle(
		"/events/{id}/attendees",
		controllers.EventsAttendees,
	).Methods("GET")

	// Have the current user attend a specific event
	router.handle(
		"/events/{id}/attend",
		router.authorizationHandler(controllers.EventsAttend),
	).Methods("POST")

	// Remove the current user for a specific event's attendees
	router.handle(
		"/events/{id}/attendance",
		router.authorizationHandler(controllers.EventsDeleteAttendance),
	).Methods("DELETE")
}
