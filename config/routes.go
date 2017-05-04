package config

import (
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/controllers"
	"github.com/omar-ozgur/flock-api/middleware"
	"github.com/urfave/negroni"
)

func InitRouter() (n *negroni.Negroni) {
	r := mux.NewRouter()

	authorizationHandler := middleware.JWTMiddleware.Handler
	r.HandleFunc("/", controllers.PagesIndex).Methods("GET")
	r.Handle("/users", authorizationHandler(controllers.UsersIndex)).Methods("GET")
	r.Handle("/users", controllers.UsersCreate).Methods("POST")
	r.Handle("/users/{id}", authorizationHandler(controllers.UsersShow)).Methods("GET")
	r.Handle("/users/{id}", authorizationHandler(controllers.UsersUpdate)).Methods("PUT")
	r.Handle("/users/{id}", authorizationHandler(controllers.UsersDelete)).Methods("DELETE")

	r.Handle("/login", controllers.UsersLogin).Methods("Post")

	r.HandleFunc("/posts", controllers.PostsIndex).Methods("GET")
	r.HandleFunc("/posts", controllers.PostsCreate).Methods("POST")
	r.HandleFunc("/posts/{id}", controllers.PostsShow).Methods("GET")
	r.HandleFunc("/posts/{id}", controllers.PostsUpdate).Methods("PUT")
	r.HandleFunc("/posts/{id}", controllers.PostsDelete).Methods("DELETE")

	r.HandleFunc("/attendees", controllers.AttendeesIndex).Methods("GET")
	r.HandleFunc("/attendees", controllers.AttendeesCreate).Methods("POST")
	r.HandleFunc("/attendees/{id}", controllers.AttendeesShow).Methods("GET")
	r.HandleFunc("/attendees/{id}", controllers.AttendeesUpdate).Methods("PUT")
	r.HandleFunc("/attendees/{id}", controllers.AttendeesDelete).Methods("DELETE")

	n = negroni.New(negroni.HandlerFunc(middleware.CustomMiddleware), negroni.NewLogger())
	n.UseHandler(r)

	return
}
