package middleware

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/omar-ozgur/flock-api/utilities"
)

// JWT middleware
var JWTMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return utilities.FLOCK_TOKEN_SECRET, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})
