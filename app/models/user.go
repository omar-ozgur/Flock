package models

import "time"

type User struct {
	Name       string
	Email      string
	Created_at time.Time
}

type Users []User
