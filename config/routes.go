package config

import (
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/controllers"
	"github.com/omar-ozgur/flock-api/middleware"
	"github.com/urfave/negroni"
)

// router holds application router information
type router struct {
	muxRouter *mux.Router
}

// Initialize the router
func (r router) init() *negroni.Negroni {

	// Create routes
	r.muxRouter = mux.NewRouter()
	r.createRoutes()

	// Initialize middleware
	negroni := middleware.Init(r.muxRouter)

	return negroni
}

func (r router) createRoutes() {

	// Create an authorization handler
	authorizationHandler := middleware.JWTMiddleware.Handler

	// Create a hanle function
	handle := r.muxRouter.Handle

	// Create user routes
	handle(
		"/signup",
		controllers.UsersCreate,
	).Methods("POST")

	handle(
		"/login",
		controllers.UsersLogin,
	).Methods("POST")

	handle(
		"/profile",
		authorizationHandler(controllers.UsersProfile),
	).Methods("Get")

	handle(
		"/users",
		controllers.UsersIndex,
	).Methods("GET")

	handle(
		"/users/search",
		controllers.UsersSearch,
	).Methods("POST")

	handle(
		"/users/{id}",
		controllers.UsersShow,
	).Methods("GET")

	handle(
		"/users/{id}",
		authorizationHandler(controllers.UsersUpdate),
	).Methods("PUT")

	handle(
		"/users/{id}",
		authorizationHandler(controllers.UsersDelete),
	).Methods("DELETE")

	handle(
		"/users/{id}/attendance",
		authorizationHandler(controllers.UsersAttendance),
	).Methods("GET")

	// Create event routes
	handle(
		"/events",
		controllers.EventsIndex,
	).Methods("GET")

	handle(
		"/events",
		authorizationHandler(controllers.EventsCreate),
	).Methods("POST")

	handle(
		"/events/search",
		controllers.EventsSearch,
	).Methods("POST")

	handle(
		"/events/{id}",
		controllers.EventsShow,
	).Methods("GET")

	handle(
		"/events/{id}",
		authorizationHandler(controllers.EventsUpdate),
	).Methods("PUT")

	handle(
		"/events/{id}",
		authorizationHandler(controllers.EventsDelete),
	).Methods("DELETE")

	handle(
		"/events/{id}/attendees",
		controllers.EventsAttendees,
	).Methods("GET")

	handle(
		"/events/{id}/attend",
		authorizationHandler(controllers.EventsAttend),
	).Methods("POST")

	handle(
		"/events/{id}/attendance",
		authorizationHandler(controllers.EventsDeleteAttendance),
	).Methods("DELETE")

	// Create login routes
	handle(
		"/loginWithFacebook",
		controllers.LoginWithFacebook,
	).Methods("POST")
}
