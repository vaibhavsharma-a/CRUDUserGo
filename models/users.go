package models

import "github.com/golang-jwt/jwt"

// Cliams contains the fields from which the Claims in JWT is created
// @Description contains jwt.StandardClaims and other to generate the claims for JWT
type Claims struct {

	// Username of the logged user
	// @example "Vaibhav sharma"
	Username string `json:"UserName"`

	// All the standard claims of the JWT token
	// @example ExpiresAt
	jwt.StandardClaims
}

// UsersInfo represents the info required to register user
// @Description takes info about user with email, username and password
type UsersInfo struct {
	// Username of the user being reigstered
	// @example "vaibhav sharma"
	Username string `json:"UserName"`

	// Userpass is the password of the user
	// @example "123@test"
	Userpass string `json:"UserPass"`

	// Email is the email address of user containing '@'
	// @example "user1@gmail.com"
	Email string `json:"EmailAddr"`
}

// UserLoginInfor represent the fields needed to login into the database
// @Description takes the info from user to log them in to the database
type UserLoginInfo struct {
	// Username is name of the user used while registring
	// @example "vaibhav sharma"
	Username string `json:"UserName"`

	// Userpass is the password that is used while registring the user
	// @example "123@test"
	Userpass string `json:"UserPass"`
}

// InfoAboutUser contains the entire info about user that is generated
// @Description contains the Id,Username,Email,Timestamp when the user created in the database
// @Id
type InfoAboutUser struct {
	// Unique identifier and is automatedly genrated at the backend
	// @example 2
	Id string `json:"UserId"`

	// Username of the logged in user
	// @example "Vaibhav sharma"
	Username string `json:"UserName"`

	// Email address of the user logged in
	// @example "test@gmail.com"
	Email string `json:"EmailAddr"`

	// The time at which the user is registered into the database
	// @eaxmple "2024-09-20 14:32:21"
	TimedCreated string `json:"CreatedAt"`
}

// UpdateUserInfo hold the info to be updated for the user in database
// @Description contains fields that can be updated by the user afer registring
type UpdateUserInfo struct {
	// Username of the logged in user
	// @example "Vaibhav sharma"
	Username string `json:"UserName,omitempty"`

	// Userpass is the password that is used while registring the user
	// @example "123@test"
	Userpass string `json:"UserPass,omitempty"`

	// Email address of the user logged in
	// @example "test@gmail.com"
	Email string `json:"EmailAddr,omitempty"`
}
