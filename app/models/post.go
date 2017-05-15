package models

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	_ "github.com/lib/pq"
	"github.com/omar-ozgur/flock-api/db"
	"github.com/omar-ozgur/flock-api/utilities"
	"gopkg.in/oleiade/reflections.v1"
	"reflect"
	"strconv"
	"time"
)

type Post struct {
	Id           int       `valid:"-"`
	Title        string    `valid:"required"`
	Description  string    `valid:"required"`
	Location     string    `valid:"required"`
	User_id      int       `valid:"required"`
	Latitude     string    `valid:"latitude,required"`
	Longitude    string    `valid:"longitude,required"`
	Zip          int       `valid:"required"`
	Time_created time.Time `valid:"-"`
	Time_expires time.Time `valid:"-"`
}

const postTableName = "posts"

var postAutoParams = map[string]bool{"Id": true, "User_id": true, "Time_created": true, "Time_expires": true}
var postRequiredParams = map[string]bool{"Title": true, "Description": true, "Location": true, "Latitude": true, "Longitude": true, "Zip": true}

func GetPosts() (status string, message string, retrievedPosts []Post) {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s;", postTableName)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), nil
	}
	rows, err := stmt.Query()
	if err != nil {
		return "error", "Failed to query posts", nil
	}

	// Get post info
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

	// Convert user ID to integer
	var err error
	post.User_id, err = strconv.Atoi(userId)
	if err != nil {
		return "error", "Invalid user ID", Post{}
	}

	// Validate user
	_, err = govalidator.ValidateStruct(post)
	if err != nil {
		return "error", fmt.Sprintf("Failed to validate post: %s", err.Error()), Post{}
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("INSERT INTO %s (", postTableName))

	// Set present column names
	var valuesStr bytes.Buffer
	var values []interface{}
	parameterIndex := 1
	queryStr.WriteString("User_id")
	valuesStr.WriteString(fmt.Sprintf("$%d", parameterIndex))
	parameterIndex += 1
	values = append(values, userId)
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldValue := fmt.Sprintf("%v", value.Field(i).Interface())
		if postAutoParams[fieldName] {
			continue
		}
		if postRequiredParams[fieldName] && reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			return "error", fmt.Sprintf("Field '%v' is not valid", fieldName), Post{}
		}
		queryStr.WriteString(fmt.Sprintf(", %v", fieldName))
		valuesStr.WriteString(fmt.Sprintf(", $%d", parameterIndex))
		values = append(values, fmt.Sprintf("%v", fieldValue))
		reflections.SetField(&post, value.Type().Field(i).Name, value.Field(i).Interface())
		parameterIndex += 1
	}

	// Finish and execute query
	queryStr.WriteString(fmt.Sprintf(") VALUES(%s", valuesStr.String()))
	queryStr.WriteString(") returning id;")
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), Post{}
	}
	err = stmt.QueryRow(values...).Scan(&post.Id)
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
	queryStr.WriteString(fmt.Sprintf("SELECT * FROM %s WHERE", postTableName))

	// Get post fields
	value := reflect.ValueOf(post)
	if value.NumField() <= 0 {
		return "error", "Invalid number of fields", nil
	}

	// Set present column names and values
	var values []interface{}
	parameterIndex := 1
	var first = true
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldValue := value.Field(i).Interface()
		if reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			continue
		}
		if !first {
			queryStr.WriteString(" AND ")
		} else {
			queryStr.WriteString(" ")
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("%v=$%d", fieldName, parameterIndex))
		values = append(values, fieldValue)
		parameterIndex += 1
	}

	// Check if any fields are valid
	if len(values) <= 0 {
		return "error", "No valid fields were found", nil
	}

	// Finish and execute query
	queryStr.WriteString(";")
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), nil
	}
	rows, err := stmt.Query(values...)
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
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", postTableName)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: %v", id)
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), Post{}
	}
	row := stmt.QueryRow(id)

	// Get post info
	var post Post
	err = row.Scan(&post.Id, &post.Title, &post.Location, &post.User_id, &post.Latitude, &post.Longitude, &post.Zip, &post.Time_created, &post.Time_expires)
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
	var values []interface{}
	parameterIndex := 1
	var first = true
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldValue := value.Field(i).Interface()
		if postAutoParams[fieldName] {
			continue
		}
		if reflect.DeepEqual(fieldValue, reflect.Zero(reflect.TypeOf(fieldValue)).Interface()) {
			continue
		}
		if !first {
			queryStr.WriteString(", ")
		} else {
			queryStr.WriteString(" ")
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("%v=$%d", fieldName, parameterIndex))
		values = append(values, fieldValue)
		parameterIndex += 1
	}

	// Finish and execute query
	queryStr.WriteString(fmt.Sprintf(" WHERE id=$%d;", parameterIndex))
	values = append(values, id)
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), Post{}
	}
	_, err = stmt.Exec(values...)
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
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", postTableName)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	utilities.Sugar.Infof("Values: %v", id)
	stmt, err := db.DB.Prepare(queryStr)
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error())
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return "error", "Failed to delete post"
	}

	return "success", "Deleted post"
}
