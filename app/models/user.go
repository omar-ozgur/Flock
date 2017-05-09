package models

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"github.com/omar-ozgur/flock-api/db"
	"github.com/omar-ozgur/flock-api/utilities"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/oleiade/reflections.v1"
	"os"
	"reflect"
	"time"
)

type User struct {
	Id           int       `valid:"-"`
	First_name   string    `valid:"alphanum,required"`
	Last_name    string    `valid:"alphanum,required"`
	Email        string    `valid:"email,required"`
	Fb_id        int       `valid:"-"`
	Password     []byte    `valid:"required"`
	Time_created time.Time `valid:"-"`
}

const userTableName = "users"

var UserAutoParams = map[string]bool{"Id": true, "Time_created": true}
var UserUniqueParams = map[string]bool{"Email": true, "Fb_id": true}
var UserRequiredParams = map[string]bool{"First_name": true, "Last_name": true, "Email": true, "Fb_id": true, "Password": true}

func CreateUser(user User) (status string, message string, createdUser User) {

	// Get user fields
	value := reflect.ValueOf(user)
	if value.NumField() <= len(UserRequiredParams) {
		return "error", "Invalid number of user parameters", User{}
	}

	// Encrypt password
	hash, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		return "error", fmt.Sprintf("Failed to encrypt password: %s", err.Error()), User{}
	}
	reflections.SetField(&user, "Password", hash)

	// Validate user
	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		return "error", fmt.Sprintf("Failed to validate user: %s", err.Error()), User{}
	}

	// Check user uniqueness
	uniqueMap := make(map[string]interface{})
	for key, _ := range UserUniqueParams {
		fieldValue, err := reflections.GetField(&user, key)
		if err != nil {
			return "error", fmt.Sprintf("Failed to get field: %s", err.Error()), User{}
		}
		uniqueMap[key] = fieldValue
	}
	status, message, retrievedUsers := SearchUsers(uniqueMap, "OR")
	if status != "success" {
		return "error", "Failed to check user uniqueness", User{}
	} else if retrievedUsers != nil {
		return "error", "User is not unique", User{}
	}

	// Create query string
	var queryStr bytes.Buffer
	queryStr.WriteString(fmt.Sprintf("INSERT INTO %s (", userTableName))

	// Set present column names
	var fieldsStr, valuesStr bytes.Buffer
	var values []interface{}
	parameterIndex := 1
	var first = true
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldValue := fmt.Sprintf("%v", value.Field(i).Interface())
		if UserAutoParams[fieldName] {
			continue
		}
		if !first {
			fieldsStr.WriteString(fmt.Sprintf(", %s", fieldName))
			valuesStr.WriteString(fmt.Sprintf(", $%d", parameterIndex))
		} else {
			fieldsStr.WriteString(fieldName)
			valuesStr.WriteString(fmt.Sprintf("$%d", parameterIndex))
			first = false
		}
		parameterIndex += 1
		values = append(values, fieldValue)
	}

	// Finish and prepare query
	queryStr.WriteString(fmt.Sprintf("%s) VALUES(%s) RETURNING id;", fieldsStr.String(), valuesStr.String()))
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	utilities.Sugar.Infof("Values: %v", values)
	stmt, err := db.DB.Prepare(queryStr.String())
	if err != nil {
		return "error", fmt.Sprintf("Failed to prepare DB query: %s", err.Error()), User{}
	}

	// Execute query
	err = stmt.QueryRow(values...).Scan(&user.Id)
	if err != nil {
		return "error", fmt.Sprintf("Failed to create new user: %s", err.Error()), User{}
	}

	// Get created user
	status, _, createdUser = GetUser(fmt.Sprintf("%v", user.Id))
	if status != "success" {
		return "error", "Failed to retrieve created user", User{}
	}

	return "success", "New user created", createdUser
}

func LoginUser(user User) (status string, message string, createdToken string) {
	if user.Email == "" || len(user.Password) == 0 {
		return "error", "Invalid login parameters", ""
	}
	var foundUser User
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE email='%s';", userTableName, user.Email)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	row := db.DB.QueryRow(queryStr)
	err := row.Scan(&foundUser.Id, &foundUser.First_name, &foundUser.Last_name, &foundUser.Email, &foundUser.Fb_id, &foundUser.Password, &foundUser.Time_created)
	if err != nil {
		return "error", "Error while retrieving user", ""
	}
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "error", "Error while encrypting password", ""
	}
	err = bcrypt.CompareHashAndPassword(hash, user.Password)
	if err != nil {
		return "error", "Error while checking password", ""
	}

	var secretKey = []byte(os.Getenv("FLOCK_TOKEN_SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = foundUser.Id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString(secretKey)
	return "success", "Login token generated", tokenString
}

func GetUser(id string) (status string, message string, retrievedUser User) {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s WHERE id=%s;", userTableName, id)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	row := db.DB.QueryRow(queryStr)

	// Get user info
	var user User
	err := row.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.Fb_id, &user.Password, &user.Time_created)
	if err != nil {
		return "error", "Failed to retrieve user information", User{}
	}

	return "success", "Retrieved user", user
}

func GetUsers() (status string, message string, retrievedUsers []User) {

	// Create and execute query
	queryStr := fmt.Sprintf("SELECT * FROM %s;", userTableName)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	rows, err := db.DB.Query(queryStr)
	if err != nil {
		return "error", "Failed to query users", nil
	}

	// Print table
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.Fb_id, &user.Password, &user.Time_created)
		if err != nil {
			return "error", "Failed to retrieve user information", nil
		}
		users = append(users, user)
	}

	return "success", "Retrieved users", users
}

func UpdateUser(id string, user User) (status string, message string, updatedUser User) {

	// Get user fields
	value := reflect.ValueOf(user)
	if value.NumField() <= 0 {
		return "error", "Invalid number of fields", user
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
			if err != nil {
				return "error", "Failed to encrypt password", User{}
			}
			queryStr.WriteString(fmt.Sprintf("%v='%v'", value.Type().Field(i).Name, hash))
		} else {
			queryStr.WriteString(fmt.Sprintf("%v='%v'", value.Type().Field(i).Name, value.Field(i).Interface()))
		}
	}

	// Finish and execute query
	queryStr.WriteString(fmt.Sprintf(" WHERE id='%s';", id))
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	_, err := db.DB.Exec(queryStr.String())
	if err != nil {
		return "error", "Failed to update user", User{}
	}

	status, message, retrievedUser := GetUser(id)
	if status == "success" {
		return "success", "Updated user", retrievedUser
	} else {
		return "error", "Failed to retrieve updated user", User{}
	}
}

func DeleteUser(id string) (status string, message string) {

	// Create and execute query
	queryStr := fmt.Sprintf("DELETE FROM %s WHERE id=%s;", userTableName, id)
	utilities.Sugar.Infof("SQL Query: %s", queryStr)
	_, err := db.DB.Exec(queryStr)
	if err != nil {
		return "error", "Failed to delete user"
	}

	return "success", "Deleted user"
}

func SearchUsers(parameters map[string]interface{}, operator string) (status string, message string, retrievedUsers []User) {

	// Create query string
	var queryStr bytes.Buffer

	// Create and execute query
	queryStr.WriteString(fmt.Sprintf("SELECT * FROM %s WHERE", userTableName))

	// Set present column names and values
	var first = true
	for key, value := range parameters {
		if !first {
			queryStr.WriteString(fmt.Sprintf(" %s ", operator))
		} else {
			queryStr.WriteString(" ")
			first = false
		}
		queryStr.WriteString(fmt.Sprintf("%v='%v'", key, value))
	}

	queryStr.WriteString(";")
	utilities.Sugar.Infof("SQL Query: %s", queryStr.String())
	rows, err := db.DB.Query(queryStr.String())
	if err != nil {
		return "error", "Failed to query users", nil
	}

	// Print table
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.Fb_id, &user.Password, &user.Time_created)
		if err != nil {
			return "error", "Failed to retrieve user information", nil
		}
		users = append(users, user)
	}

	return "success", "Retrieved users", users
}
