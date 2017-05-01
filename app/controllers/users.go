package controllers

import (
	"encoding/json"
	"github.com/omar-ozgur/flock-api/app/models"
	"net/http"
)

func UsersIndex(w http.ResponseWriter, r *http.Request) {
	user := models.User{First_name: "Omar", Last_name: "Ozgur", Email: "oozgur217@gmail.com", Fb_id: 1}
	models.SaveUser(user)
	models.QueryUsers()

	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}
