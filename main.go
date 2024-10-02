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

	_ "dbgo/docs"

	_ "github.com/go-sql-driver/mysql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var db *sql.DB

func main() {

	// @title						This is a RESTFul CRUD API
	// @version					1.0
	// @description			This API allows to user to create,get info about,update & delete a user

	// @contact.name		vaibhav sharma
	// @contact.url			https://www.linkedin.com/in/sharmaaavaibhav/
	// @contact.email		vaibhav1863sharma@gmail.com

	// @license.name 		MIT
	// @license.url			https://opensource.org/license/mit

	// @host 						localhost:8080
	// @basepath				/

	//? loading the env file if the program runs locally
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println(err)
		}
	}

	db_user := os.Getenv("DB_User")
	db_passwrd := os.Getenv("DB_Pass")
	db_host := os.Getenv("DB_Host")
	db_port := os.Getenv("DB_Port")
	db_dbname := os.Getenv("DB_Dbname")

	//? creating a database connection
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

	//? deffering the closer of the database connection once the main funciton is done running
	defer db.Close()

	//? creating a gin router
	router := gin.Default()

	//? user registing endpoint
	router.POST("/register", func(c *gin.Context) {
		handlers.RegisterUserHandler(db, c)
	})

	//? user login endpoint
	router.POST("/login", func(c *gin.Context) {
		handlers.LoginUserHandler(db, c)
	})

	//? middelware protected user info endpoint
	router.GET("/user/:username", middleware.JWTAuthMiddlerware(), func(c *gin.Context) {
		handlers.GetUserByUsernameHandler(db, c)
	})

	//? middelware protected user deletion endpoint
	router.DELETE("deluser/:username", middleware.JWTAuthMiddlerware(), func(c *gin.Context) {
		handlers.DeleteUserByIdHanlder(db, c)
	})

	//? middleware protected user updation endpoint
	router.PUT("/update/:username", middleware.JWTAuthMiddlerware(), func(c *gin.Context) {
		handlers.UpdateUserByUsernameHandler(db, c)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback to default port for local development
	}

	router.Run(":" + port)

}
