package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/models"
	"io/ioutil"
	"net/http"
)

func UsersIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users := models.GetUsers()
	j, _ := json.Marshal(users)
	w.Write(j)
	fmt.Println("Retrieved users")
}

func UsersShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	user := models.GetUser(vars["id"])
	j, _ := json.Marshal(user)
	w.Write(j)
	fmt.Println("Retrieved user")
}

func UsersCreate(w http.ResponseWriter, r *http.Request) {
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
}

func UsersUpdate(w http.ResponseWriter, r *http.Request) {
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
}

func UsersDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	models.DeleteUser(vars["id"])
	fmt.Println("Deleted user")
}
