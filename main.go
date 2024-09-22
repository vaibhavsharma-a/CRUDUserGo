package main

import (
	"database/sql"
	"dbgo/handlers"
	"dbgo/middleware"
	"fmt"
	"log"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}
	db_user := os.Getenv("DB_user")
	db_passwrd := os.Getenv("DB_Pass")
	db_host := os.Getenv("DB_Host")
	db_port := os.Getenv("DB_Port")
	db_dbname := os.Getenv("DB_Dbname")

	var err error
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_user, db_passwrd, db_host, db_port, db_dbname)

	db, err = sql.Open("mysql", dns)
	if err != nil {
		log.Fatal("Error while opening database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("can't ping to database", err)
	}
	fmt.Println("Connected to database")
	if db == nil {
		log.Println("Database connection is NIL")
	}

	defer db.Close()

	router := gin.Default()

	router.POST("/register", func(c *gin.Context) {
		handlers.RegisterUserHandler(db, c)
	})

	router.POST("/login", func(c *gin.Context) {
		handlers.LoginUserHandler(db, c)
	})

	router.GET("/user/:username", middleware.JWTAuthMiddlerware(), func(c *gin.Context) {
		handlers.GetUserByUsernameHandler(db, c)
	})

	router.DELETE("deluser/:username", middleware.JWTAuthMiddlerware(), func(c *gin.Context) {
		handlers.DeleteUserByIdHanlder(db, c)
	})

	router.PUT("/update/:username", middleware.JWTAuthMiddlerware(), func(c *gin.Context) {
		handlers.UpdateUserByUsernameHandler(db, c)
	})
	router.Run(":8080")

}
