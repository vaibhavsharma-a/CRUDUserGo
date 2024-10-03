package handlers

import (
	"database/sql"
	"dbgo/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserHandler  			godoc
// @Summary										Register user to the Database
// @Description								Take user info and update it to the database
// @tags											userinfo
// @Accept										json
// @Produce										json
// @Param											userinfo body models.UsersInfo true "Info about the user"
// @Success										201 {string} string "The user {username} added successfully in database"
// @Failure										400 {object} map[string]string "Error: Bad Request"
// @Failure										500 {object} map[string]string "Error: Failed to hash the password"
// @Failure										409 {object} map[string]string "Error: Username {username} already exists in the database"
// @Router                    /register [post]
func RegisterUserHandler(db *sql.DB, c *gin.Context) {
	var user models.UsersInfo

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	//!Encrypt the password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Userpass), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to hash the password"})
		return
	}
	//!Validate the entered values
	var count int

	if err := db.QueryRow("SELECT COUNT(*) FROM users WHERE UserName = ?", user.Username).Scan(&count); err != nil {
		log.Println("Validation error", err)
		return
	}
	if count == 1 {
		c.JSON(http.StatusConflict, gin.H{"Error": fmt.Sprintf("Username %s already exists in the database", user.Username)})
		return
	}

	insertStmt, err := db.Prepare("INSERT INTO users (UserName,UserPass,EmailAddr) VALUES (?,?,?)")
	if err != nil {
		log.Printf("Error preparing Sql statement: %s", err)
		return
	}
	defer insertStmt.Close()

	if _, err := insertStmt.Exec(user.Username, hashedPass, user.Email); err != nil {
		log.Printf("Error Executing insert statement: %s", err)
		return

	}
	c.JSON(http.StatusAccepted, gin.H{"Message": fmt.Sprintf("The user %s added successfully to the database", user.Username)})

}

// LoginUserHandler  			godoc
// @Summary								Login the registered users and generate JWT token
// @Description						Authenticate user and return a JWT token
// @tags									Authentication
// @Accept								json
// @Produce								json
// @Param									userlogininfo body models.UserLoginInfo true "User Login Credentials"
// @Success								201 {object} map[string]string "Message : User is successfully logged in!, token : tokenstring"
// @Failure								400 {object} map[string]string "Error: Invalid Credentials"
// @Failure								401 {object} map[string]string "Error: UserName {username} is not present in the database"
// @Failure								500 {object} map[string]string "Error: There is some interanl sever error while fetching"
// @Failure								400 {object} map[string]string "Error: Could not retrieve the password from the database"
// @Failure								401 {object} map[string]string "Error: invalid passowrd"
// @Failure								500 {object} map[string]string "Error: Could not sign the token with secret"
// @Router                /login [post]
func LoginUserHandler(db *sql.DB, c *gin.Context) {
	var userlogin models.UserLoginInfo

	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println(err)
		}
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	if err := c.ShouldBindJSON(&userlogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Credentials"})
		return
	}

	var usrnm string

	//? Query the database for the username
	err := db.QueryRow("SELECT UserName from users WHERE UserName = ?", userlogin.Username).Scan(&usrnm)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": fmt.Sprintf("UserName %s is not present in the database", userlogin.Username)})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "There is some internal Error in the server"})
		return
	}

	var storedHashPass string

	if err := db.QueryRow("SELECT UserPass FROM users WHERE UserName = ?", usrnm).Scan(&storedHashPass); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Could not retrieve the Password from Database"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHashPass), []byte(userlogin.Userpass)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid Password"})
		return
	}
	//!JWT TOKEN Generation

	claims := models.Claims{
		Username: userlogin.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//! signing the token with signing string
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Could not sing the token with secret"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User is successfully logged in!",
		"token":   tokenString,
	})

}

// GetUserByNameHandler  			godoc
// @Summary										Get details of logged in user
// @Description								This Routes take a JWT token for authentication and retrieves details of the logged in user
// @tags											users
// @Accept										json
// @Produce										json
// @Param											Authorization header string true "JWT Token"
// @Success										200 {object} models.InfoAboutUser "User info retrieved successfully"
// @Failure										404 {object} map[string]string "Error: There is no such user in the database"
// @Failure										500 {object} map[string]string "Error: Internal Server error"
// @Router                    /user/:username [get]
// @Security									BearerAuth
func GetUserByUsernameHandler(db *sql.DB, c *gin.Context) {
	username := c.MustGet("UserName").(string)

	var info models.InfoAboutUser

	err := db.QueryRow("SELECT UserId,UserName,EmailAddr,CreatedAt From users WHERE UserName = ?", username).Scan(&info.Id, &info.Username, &info.Email, &info.TimedCreated)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"Error": "There is no such user in the database"})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Internal server error"})
	}

	infromation := models.InfoAboutUser{
		Id:           info.Id,
		Username:     info.Username,
		Email:        info.Email,
		TimedCreated: info.TimedCreated,
	}
	c.JSON(http.StatusOK, gin.H{"Message": "User info retrieved successfully",
		"information": infromation,
	})

}

// DeleteUserByNameHandler  	godoc
// @Summary										Delete the logged in user
// @Description								This Routes take a JWT token for authentication and deletes the logged in user
// @tags											users
// @Accept										json
// @Produce										json
// @Param											Authorization header string true "JWT Token"
// @Success										200 {object} map[string]string "{username} has been deleted from the database, You may logout"
// @Failure										401 {object} map[string]string "Error: {username} does not exsits in the database"
// @Failure										500 {object} map[string]string "Error: There is some Internal Server error"
// @Router                    /deluser/:username [delete]
// @Security									BearerAuth
func DeleteUserByNameHandler(db *sql.DB, c *gin.Context) {
	username := c.MustGet("UserName").(string)

	deleteStmt, err := db.Prepare("DELETE FROM users WHERE UserName = ?")

	if err != nil {
		log.Printf("Error while preparing Query: %s", err)
		return
	}
	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(username)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": fmt.Sprintf("%s Does not exist in the database", username)})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "There is some internal server error"})
	}

	c.JSON(http.StatusOK, gin.H{"Deleted": fmt.Sprintf("%s has been deleted from the database, You may logout", username)})

}

// UpdateUserByNameHandler  	godoc
// @Summary										Update the information of logged in user
// @Description								Takes the information to be updated and required JWT token for the authorization
// @tags											users
// @Accept										json
// @Produce										json
// @Param											Authorization header string true "JWT Token"
// @Param											updateInfo 		body 	 models.UpdateUserInfo true "Info about the user to be updated"
// @Success										200 {object} map[string]string "{username} Info is successfully updated"
// @Failure										400 {object} map[string]string "Error: Invalid Inputs"
// @Failure										500 {object} map[string]string "Error: Failed to hash the password"
// @Failure										500 {object} map[string]string "Error: Failed to Update info"
// @Router                    /update/:username [update]
// @Security									BearerAuth
func UpdateUserByNameHandler(db *sql.DB, c *gin.Context) {
	username := c.MustGet("UserName").(string)

	var updtInfo models.UpdateUserInfo

	if err := c.ShouldBindJSON(&updtInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Inputs"})
		return
	}

	if updtInfo.Userpass != "" {
		hashUpdatePass, err := bcrypt.GenerateFromPassword([]byte(updtInfo.Userpass), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to hash the password"})
			return
		}
		updtInfo.Userpass = string(hashUpdatePass)
	}

	_, err := db.Exec(`
				UPDATE users
				SET
					UserName = COALESCE(NULLIF(?,''),UserName),
					EmailAddr = COALESCE(NULLIF(?,''),EmailAddr),
					UserPass = COALESCE(NULLIF(?,''),UserPass)
				WHERE UserName = ?`, updtInfo.Username, updtInfo.Email, updtInfo.Userpass, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to Update the Info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Success": fmt.Sprintf("%s Info is successfully Updated", username)})

}
