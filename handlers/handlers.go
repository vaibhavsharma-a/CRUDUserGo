package handlers

import (
	"database/sql"
	"dbgo/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

}

func GetAllUsersHandler(db *sql.DB, c *gin.Context) {

}

func GetUserByIdHandler(db *sql.DB, c *gin.Context) {

}

func DeleteUserByIdHanlder(db *sql.DB, c *gin.Context) {

}
