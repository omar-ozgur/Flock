package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/omar-ozgur/flock-api/app/models"
	"io/ioutil"
	"net/http"
)

var PostsIndex = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	posts := models.GetPosts()
	j, _ := json.Marshal(posts)
	w.Write(j)
	fmt.Println("Retrieved posts")
})

var PostsShow = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	post := models.GetPost(vars["id"])
	j, _ := json.Marshal(post)
	w.Write(j)
	fmt.Println("Retrieved post")
})

var PostsCreate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post models.Post
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &post)
	status := models.CreatePost(post)
	if status == true {
		fmt.Println("Created new post")
	} else {
		fmt.Println("New post is not valid")
	}
})

var PostsUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	var post models.Post
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &post)
	status := models.UpdatePost(vars["id"], post)
	if status == true {
		fmt.Println("Updated post")
	} else {
		fmt.Println("Post info is not valid")
	}
})

var PostsDelete = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	models.DeletePost(vars["id"])
	fmt.Println("Deleted post")
})
