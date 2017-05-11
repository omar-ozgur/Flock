package config

import (
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/controllers"
	"github.com/omar-ozgur/flock-api/middleware"
	"github.com/urfave/negroni"
)

func InitRouter() (n *negroni.Negroni) {
	authorizationHandler := middleware.JWTMiddleware.Handler

	r := mux.NewRouter()

	r.HandleFunc("/", controllers.PagesIndex).Methods("GET")

	r.Handle("/signup", controllers.UsersCreate).Methods("POST")
	r.Handle("/login", controllers.UsersLogin).Methods("POST")
	r.Handle("/profile", authorizationHandler(controllers.UsersProfile)).Methods("Get")
	r.Handle("/users", controllers.UsersIndex).Methods("GET")
	r.Handle("/users/{id}", controllers.UsersShow).Methods("GET")
	r.Handle("/users/{id}", authorizationHandler(controllers.UsersUpdate)).Methods("PUT")
	r.Handle("/users/{id}", authorizationHandler(controllers.UsersDelete)).Methods("DELETE")

	r.Handle("/posts", controllers.PostsIndex).Methods("GET")
	r.Handle("/posts", authorizationHandler(controllers.PostsCreate)).Methods("POST")
	r.Handle("/posts/search", controllers.PostsSearch).Methods("POST")
	r.Handle("/posts/{id}", controllers.PostsShow).Methods("GET")
	r.Handle("/posts/{id}", authorizationHandler(controllers.PostsUpdate)).Methods("PUT")
	r.Handle("/posts/{id}", authorizationHandler(controllers.PostsDelete)).Methods("DELETE")

	r.Handle("/posts/{id}/attendees", controllers.PostsAttendees).Methods("GET")
	r.Handle("/posts/{id}/attend", authorizationHandler(controllers.PostsAttend)).Methods("POST")
	r.Handle("/posts/{id}/attendance", authorizationHandler(controllers.PostsDeleteAttendance)).Methods("DELETE")

	//login handling
	r.HandleFunc("/login", controllers.LoginIndex).Methods("GET")

	n = negroni.New(negroni.HandlerFunc(middleware.CustomMiddleware), negroni.NewLogger())
	n.UseHandler(r)

	return
}
