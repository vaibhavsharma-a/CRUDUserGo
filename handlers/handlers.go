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

func LoginUserHandler(db *sql.DB, c *gin.Context) {
	var userlogin models.UserLoginInfo
	if err := godotenv.Load(); err != nil {
		log.Println(err)
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

/*func GetAllUsersHandler(db *sql.DB, c *gin.Context) {

}*/

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
	c.JSON(http.StatusOK, infromation)

}

func DeleteUserByIdHanlder(db *sql.DB, c *gin.Context) {

}
