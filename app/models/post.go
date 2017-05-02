package models

import (
	"bytes"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/omar-ozgur/flock-api/db"
	"github.com/omar-ozgur/flock-api/utilities"
	"reflect"
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

var postAutoParams = map[string]bool{"Id": true, "Time_created": true, "Time_expires": true}
var postRequiredParams = map[string]bool{"Title": true, "Location": true, "User_id": true, "Latitude": true, "Longitude": true, "Zip": true}

func CreatePost(post Post) bool {

	// Get post fields
	value := reflect.ValueOf(post)
	if value.NumField() <= len(postRequiredParams) {
		return false
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("INSERT INTO %s (", postTableName))

	// Set present column names
	var first = true
	var values []string
	for i := 0; i < value.NumField(); i++ {
		name := value.Type().Field(i).Name
		if postAutoParams[name] {
			continue
		}
		if postRequiredParams[name] && reflect.DeepEqual(value.Field(i).Interface(), reflect.Zero(reflect.TypeOf(value.Field(i).Interface())).Interface()) {
			return false
		}
		if !first {
			queryStr.WriteString(", ")
		} else {
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("%v", value.Type().Field(i).Name))
		values = append(values, fmt.Sprintf("%v", value.Field(i).Interface()))
	}

	// Set present column values
	queryStr.WriteString(") VALUES(")
	first = true
	for i := 0; i < len(values); i++ {
		if !first {
			queryStr.WriteString(", ")
		} else {
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("'%v'", values[i]))
	}

	// Finish and execute query
	queryStr.WriteString(")")
	fmt.Println("SQL Query:", queryStr.String())
	_, err := db.DB.Exec(queryStr.String())
	utilities.CheckErr(err)

	return true
}

func UpdatePost(id string, post Post) bool {

	// Get post fields
	value := reflect.ValueOf(post)
	if value.NumField() <= 0 {
		return false
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
	queryStr.WriteString(fmt.Sprintf(" WHERE id='%s'", id))
	fmt.Println("SQL Query:", queryStr.String())
	_, err := db.DB.Exec(queryStr.String())
	utilities.CheckErr(err)

	return true
}

func DeletePost(id string) {

	// Create and execute query
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE id=%s", postTableName, id)
	fmt.Println("SQL Query:", queryStr)
	_, err := db.DB.Exec(queryStr)
	utilities.CheckErr(err)
}

func GetPost(id string) Post {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE id=%s", postTableName, id)
	fmt.Println("SQL Query:", queryStr)
	row := db.DB.QueryRow(queryStr)

	// Get post info
	var post Post
	err := row.Scan(&post.Id, &post.Title, &post.Location, &post.User_id, &post.Latitude, &post.Longitude, &post.Zip, &post.Time_created, &post.Time_expires)
	utilities.CheckErr(err)

	return post
}

func GetPosts() []Post {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s", postTableName)
	fmt.Println("SQL Query:", queryStr)
	rows, err := db.DB.Query(queryStr)
	utilities.CheckErr(err)

	// Print table
	var posts []Post
	fmt.Printf(" %-5v | %-20v | %-20v | %-5v | %-20v | %-20v | %-5v | %-20v | %-20v\n", "id", "title", "location", "user_id", "latitude", "longitude", "zip", "time_created", "time_expires")
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Title, &post.Location, &post.User_id, &post.Latitude, &post.Longitude, &post.Zip, &post.Time_created, &post.Time_expires)
		utilities.CheckErr(err)
		posts = append(posts, post)
		fmt.Printf(" %-5v | %-20v | %-20v | %-5v | %-20v | %-20v | %-5v | %-20v | %-20v\n", post.Id, post.Title, post.Location, post.User_id, post.Latitude, post.Longitude, post.Zip, post.Time_created, post.Time_expires)
	}

	return posts
}
