package models

import "github.com/golang-jwt/jwt"

type Claims struct {
	Username string `json:"UserName"`
	jwt.StandardClaims
}

type UsersInfo struct {
	Username string `json:"UserName"`
	Userpass string `json:"UserPass"`
	Email    string `json:"EmailAddr"`
}

type UserLoginInfo struct {
	Username string `json:"UserName"`
	Userpass string `json:"UserPass"`
}
