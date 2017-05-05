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

	users := models.GetUsers()
	j, _ := json.Marshal(users)
	w.Write(j)
	fmt.Println("Retrieved users")
})

var UsersProfile = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	user := models.GetUser(current_user_id)
	j, _ := json.Marshal(user)
	w.Write(j)
	fmt.Println("Retrieved current user")
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
	j, _ := json.Marshal(user)
	w.Write(j)
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
	foundUser := models.GetUser(vars["id"])
	j, _ := json.Marshal(foundUser)
	w.Write(j)
	if status == true {
		fmt.Println("Updated user")
	} else {
		fmt.Println("User info is not valid")
	}
})

var UsersDelete = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	models.DeleteUser(vars["id"])
	j, _ := json.Marshal("Deleted user")
	w.Write(j)
	fmt.Println("Deleted user")
})

var UsersLogin = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var user models.User
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &user)
	loginToken := models.LoginUser(user)
	j, _ := json.Marshal(loginToken)
	w.Write(j)
	if loginToken != "" {
		fmt.Println("Retrieved login token: %v", loginToken)
	} else {
		fmt.Println("Could not log in")
	}
})
