package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/omar-ozgur/flock-api/app/models"
	"io/ioutil"
	"net/http"
)

func UsersIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users := models.QueryUsers()
	j, _ := json.Marshal(users)
	w.Write(j)
}

func UsersCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &user)
	fmt.Println(user)
	status := models.SaveUser(user)
	if status == true {
		fmt.Println("Inserted new user")
	} else {
		fmt.Println("New user is not valid")
	}
}
