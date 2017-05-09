package models

import (
	"bytes"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/omar-ozgur/flock-api/db"
	"github.com/omar-ozgur/flock-api/utilities"
	"gopkg.in/oleiade/reflections.v1"
	"reflect"
	"strconv"
	"time"
)

type Post struct {
	Id           int
	Title        string
	Location     string
	User_id      int
	Latitude     float64
	Longitude    float64
	Zip          int
	Time_created time.Time
	Time_expires time.Time
}

const postTableName = "posts"

var postAutoParams = map[string]bool{"Id": true, "User_id": true, "Time_created": true, "Time_expires": true}
var postRequiredParams = map[string]bool{"Title": true, "Location": true, "Latitude": true, "Longitude": true, "Zip": true}

func GetPosts() (status string, message string, retrievedPosts []Post) {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s;", postTableName)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	rows, err := db.DB.Query(queryStr)
	if err != nil {
		return "error", "Failed to query posts", nil
	}

	// Print table
	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Title, &post.Location, &post.User_id, &post.Latitude, &post.Longitude, &post.Zip, &post.Time_created, &post.Time_expires)
		if err != nil {
			return "error", "Failed to retrieve post information", nil
		}
		posts = append(posts, post)
	}

	return "success", "Retrieved posts", posts
}

func CreatePost(userId string, post Post) (status string, message string, createdPost Post) {

	// Get post fields
	value := reflect.ValueOf(post)
	if value.NumField() <= len(postRequiredParams) {
		return "error", "Invalid post parameters", Post{}
	}

	var err error
	post.User_id, err = strconv.Atoi(userId)
	if err != nil {
		return "error", "Invalid user ID", Post{}
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("INSERT INTO %s (", postTableName))

	// Set present column names
	var values []string
	queryStr.WriteString("User_id")
	values = append(values, userId)
	for i := 0; i < value.NumField(); i++ {
		name := value.Type().Field(i).Name
		if postAutoParams[name] {
			continue
		}
		if postRequiredParams[name] && reflect.DeepEqual(value.Field(i).Interface(), reflect.Zero(reflect.TypeOf(value.Field(i).Interface())).Interface()) {
			return "error", fmt.Sprintf("Field '%v' is not valid", value.Type().Field(i).Name), Post{}
		}
		queryStr.WriteString(", ")
		queryStr.WriteString(fmt.Sprintf("%v", value.Type().Field(i).Name))
		values = append(values, fmt.Sprintf("%v", value.Field(i).Interface()))
		reflections.SetField(&post, value.Type().Field(i).Name, value.Field(i).Interface())
	}

	// Set present column values
	queryStr.WriteString(") VALUES(")
	first := true
	for i := 0; i < len(values); i++ {
		if !first {
			queryStr.WriteString(", ")
		} else {
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("'%v'", values[i]))
	}

	// Finish and execute query
	queryStr.WriteString(") returning id;")
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	err = db.DB.QueryRow(queryStr.String()).Scan(&post.Id)
	if err != nil {
		return "error", "Failed to create new post", Post{}
	}

	// Create attendee
	attendee := Attendee{Post_id: post.Id, User_id: post.User_id}
	fmt.Println(attendee)
	status, _, _ = CreateAttendee(attendee)
	if status == "success" {
		return "success", "New post created", post
	} else {
		return "error", "Failed to add attendee to new post", Post{}
	}
}

func SearchPosts(post Post) (status string, message string, retrievedPosts []Post) {

	// Create query string
	var queryStr bytes.Buffer

	// Create and execute query
	queryStr.WriteString(fmt.Sprintf("SELECT * FROM %s WHERE", postTableName))

	// Get post fields
	value := reflect.ValueOf(post)
	if value.NumField() <= 0 {
		return "error", "Invalid number of fields", nil
	}

	// Set present column names and values
	var first = true
	for i := 0; i < value.NumField(); i++ {
		if reflect.DeepEqual(value.Field(i).Interface(), reflect.Zero(reflect.TypeOf(value.Field(i).Interface())).Interface()) {
			continue
		}
		if !first {
			queryStr.WriteString(" AND ")
		} else {
			queryStr.WriteString(" ")
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("%v='%v'", value.Type().Field(i).Name, value.Field(i).Interface()))
	}

	queryStr.WriteString(";")
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	rows, err := db.DB.Query(queryStr.String())
	if err != nil {
		return "error", "Failed to query posts", nil
	}

	// Print table
	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Title, &post.Location, &post.User_id, &post.Latitude, &post.Longitude, &post.Zip, &post.Time_created, &post.Time_expires)
		if err != nil {
			return "error", "Failed to retrieve post information", nil
		}
		posts = append(posts, post)
	}

	return "success", "Retrieved posts", posts
}

func GetPost(id string) (status string, message string, retrievedPost Post) {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE id=%s;", postTableName, id)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	row := db.DB.QueryRow(queryStr)

	// Get post info
	var post Post
	err := row.Scan(&post.Id, &post.Title, &post.Location, &post.User_id, &post.Latitude, &post.Longitude, &post.Zip, &post.Time_created, &post.Time_expires)
	if err != nil {
		return "error", "Failed to retrieve post information", Post{}
	}

	return "success", "Retrieved post", post
}

func UpdatePost(id string, post Post) (status string, message string, updatedPost Post) {

	// Get post fields
	value := reflect.ValueOf(post)
	if value.NumField() <= 0 {
		return "error", "Invalid number of fields", Post{}
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("UPDATE %s SET", postTableName))

	// Set present column names and values
	var first = true
	for i := 0; i < value.NumField(); i++ {
		name := value.Type().Field(i).Name
		if postAutoParams[name] {
			continue
		}
		if reflect.DeepEqual(value.Field(i).Interface(), reflect.Zero(reflect.TypeOf(value.Field(i).Interface())).Interface()) {
			continue
		}
		if !first {
			queryStr.WriteString(", ")
		} else {
			queryStr.WriteString(" ")
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("%v='%v'", value.Type().Field(i).Name, value.Field(i).Interface()))
	}

	// Finish and execute query
	queryStr.WriteString(fmt.Sprintf(" WHERE id='%s';", id))
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	_, err := db.DB.Exec(queryStr.String())
	if err != nil {
		return "error", "Failed to update post", Post{}
	}

	status, message, retrievedPost := GetPost(id)
	if status == "success" {
		return "success", "Updated post", retrievedPost
	} else {
		return "error", "Failed to retrieve updated post", Post{}
	}
}

func DeletePost(id string) (status string, message string) {

	// Create and execute query
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE id=%s;", postTableName, id)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	_, err := db.DB.Exec(queryStr)
	if err != nil {
		return "error", "Failed to delete post"
	}

	return "success", "Deleted post"
}
