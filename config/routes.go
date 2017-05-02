package config

import (
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/controllers"
	"github.com/omar-ozgur/flock-api/middleware"
	"github.com/urfave/negroni"
)

func InitRouter() (n *negroni.Negroni) {
	r := mux.NewRouter()

	r.HandleFunc("/", controllers.PagesIndex).Methods("GET")
	r.HandleFunc("/users", controllers.UsersIndex).Methods("GET")
	r.HandleFunc("/users", controllers.UsersCreate).Methods("POST")
	r.HandleFunc("/users/{id}", controllers.UsersUpdate).Methods("PUT")
	r.HandleFunc("/users/{id}", controllers.UsersDelete).Methods("DELETE")

	n = negroni.New(negroni.HandlerFunc(middleware.CustomMiddleware), negroni.NewLogger())
	n.UseHandler(r)

	return
}
