package models

import (
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"github.com/omar-ozgur/flock-api/db"
	"github.com/omar-ozgur/flock-api/utilities"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"time"
)

type User struct {
	Id           int
	First_name   string
	Last_name    string
	Email        string
	Fb_id        int
	Password     []byte
	Time_created time.Time
}

const userTableName = "users"

var UserAutoParams = map[string]bool{"Id": true, "Time_created": true}
var UserRequiredParams = map[string]bool{"First_name": true, "Last_name": true, "Email": true, "Fb_id": true, "Password": true}

func CreateUser(user User) bool {

	// Get user fields
	value := reflect.ValueOf(user)
	if value.NumField() <= len(UserRequiredParams) {
		return false
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("INSERT INTO %s (", userTableName))

	// Set present column names
	var first = true
	var values []string
	for i := 0; i < value.NumField(); i++ {
		name := value.Type().Field(i).Name
		if UserAutoParams[name] {
			continue
		}
		if UserRequiredParams[name] && reflect.DeepEqual(value.Field(i).Interface(), reflect.Zero(reflect.TypeOf(value.Field(i).Interface())).Interface()) {
			return false
		}
		if !first {
			queryStr.WriteString(", ")
		} else {
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("%v", value.Type().Field(i).Name))
		if name == "Password" {
			hash, err := bcrypt.GenerateFromPassword(value.Field(i).Interface().([]byte), bcrypt.DefaultCost)
			utilities.CheckErr(err)
			fmt.Println("Hash to store:", string(hash))
			values = append(values, fmt.Sprintf("%v", hash))
		} else {
			values = append(values, fmt.Sprintf("%v", value.Field(i).Interface()))
		}
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

func UpdateUser(id string, user User) bool {

	// Get user fields
	value := reflect.ValueOf(user)
	if value.NumField() <= 0 {
		return false
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("UPDATE %s SET", userTableName))

	// Set present column names and values
	var first = true
	for i := 0; i < value.NumField(); i++ {
		name := value.Type().Field(i).Name
		if UserAutoParams[name] {
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
		if name == "Password" {
			hash, err := bcrypt.GenerateFromPassword([]byte(value.Field(i).Interface().(string)), bcrypt.DefaultCost)
			utilities.CheckErr(err)
			fmt.Println("Hash to store:", string(hash))
			queryStr.WriteString(fmt.Sprintf("%v='%v'", value.Type().Field(i).Name, hash))
		} else {
			queryStr.WriteString(fmt.Sprintf("%v='%v'", value.Type().Field(i).Name, value.Field(i).Interface()))
		}
	}

	// Finish and execute query
	queryStr.WriteString(fmt.Sprintf(" WHERE id='%s'", id))
	fmt.Println("SQL Query:", queryStr.String())
	_, err := db.DB.Exec(queryStr.String())
	utilities.CheckErr(err)

	return true
}

func DeleteUser(id string) {

	// Create and execute query
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE id=%s", userTableName, id)
	fmt.Println("SQL Query:", queryStr)
	_, err := db.DB.Exec(queryStr)
	utilities.CheckErr(err)
}

func GetUser(id string) User {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE id=%s", userTableName, id)
	fmt.Println("SQL Query:", queryStr)
	row := db.DB.QueryRow(queryStr)

	// Get user info
	var user User
	err := row.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.Fb_id, &user.Password, &user.Time_created)
	utilities.CheckErr(err)

	return user
}

func GetUsers() []User {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s", userTableName)
	fmt.Println("SQL Query:", queryStr)
	rows, err := db.DB.Query(queryStr)
	utilities.CheckErr(err)

	// Print table
	var users []User
	fmt.Printf(" %-5v | %-20v | %-20v | %-20v | %-20v | %-20v | %-20v\n", "id", "first_name", "last_name", "email", "fb_id", "password", "time_created")
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.Fb_id, &user.Password, &user.Time_created)
		utilities.CheckErr(err)
		users = append(users, user)
		fmt.Printf(" %-5v | %-20v | %-20v | %-20v | %-20v | %-20v | %-20v\n", user.Id, user.First_name, user.Last_name, user.Email, user.Fb_id, user.Password, user.Time_created)
	}

	return users
}

func LoginUser(user User) string {
	if user.Email == "" || len(user.Password) == 0 {
		return ""
	}
	var foundUser User
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE email='%s'", userTableName, user.Email)
	fmt.Println("SQL Query:", queryStr)
	row := db.DB.QueryRow(queryStr)
	err := row.Scan(&foundUser.Id, &foundUser.First_name, &foundUser.Last_name, &foundUser.Email, &foundUser.Fb_id, &foundUser.Password, &foundUser.Time_created)
	if err != nil {
		return ""
	}
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	err = bcrypt.CompareHashAndPassword(hash, user.Password)
	if err != nil {
		return ""
	}

	var secretKey = []byte("secret")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["first_name"] = user.First_name
	claims["last_name"] = user.Last_name
	claims["email"] = user.email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString(secretKey)
	return tokenString
}
