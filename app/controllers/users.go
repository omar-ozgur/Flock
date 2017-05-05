package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/models"
	"github.com/omar-ozgur/flock-api/utilities"
	"io/ioutil"
	"net/http"
)

var UsersIndex = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := utilities.GetClaims(r.Header.Get("Authorization"))
	fmt.Println(claims["user_id"])

	users := models.GetUsers()
	j, _ := json.Marshal(users)
	w.Write(j)
	fmt.Println("Retrieved users")
})

var UsersShow = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	user := models.GetUser(vars["id"])
	j, _ := json.Marshal(user)
	w.Write(j)
	fmt.Println("Retrieved user")
})

var UsersCreate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &user)
	status := models.CreateUser(user)
	if status == true {
		fmt.Println("Created new user")
	} else {
		fmt.Println("New user is not valid")
	}
})

var UsersUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	var user models.User
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &user)
	status := models.UpdateUser(vars["id"], user)
	if status == true {
		fmt.Println("Updated user")
	} else {
		fmt.Println("User info is not valid")
	}
})

var UsersDelete = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	models.DeleteUser(vars["id"])
	fmt.Println("Deleted user")
})

var UsersLogin = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var user models.User
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &user)
	loginToken := models.LoginUser(user)
	if loginToken != "" {
		fmt.Println(loginToken)
	} else {
		fmt.Println("Could not log in")
	}
})
