package controllers

import (
	"encoding/json"
	"github.com/omar-ozgur/flock-api/app/models"
	"net/http"
)

func UsersIndex(w http.ResponseWriter, r *http.Request) {
	users := models.Users{
		models.User{Name: "Omar", Email: "oozgur217@gmail.com"},
		models.User{Name: "Aditya", Email: "rajuaditya@gmail.com"},
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		panic(err)
	}
}
