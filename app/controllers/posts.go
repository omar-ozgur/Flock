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

var PostsIndex = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status, message, retrievedPosts := models.GetPosts()

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"posts":   retrievedPosts,
	})
	w.Write(JSON)
})

var PostsCreate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := utilities.GetClaims(r.Header.Get("Authorization")[len("Bearer "):])
	current_user_id := fmt.Sprintf("%v", claims["user_id"])

	var post models.Post
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &post)

	status, message, createdPost := models.CreatePost(current_user_id, post)

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"post":    createdPost,
	})
	w.Write(JSON)
})

var PostsShow = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	status, message, retrievedPost := models.GetPost(vars["id"])

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"post":    retrievedPost,
	})
	w.Write(JSON)
})

var PostsUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var post models.Post
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &post)

	status, message, updatedPost := models.UpdatePost(vars["id"], post)

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
		"post":    updatedPost,
	})
	w.Write(JSON)
})

var PostsDelete = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	status, message := models.DeletePost(vars["id"])

	JSON, _ := json.Marshal(map[string]interface{}{
		"status":  status,
		"message": message,
	})
	w.Write(JSON)
})
