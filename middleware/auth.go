package middleware

import (
	"dbgo/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func JWTAuthMiddlerware() gin.HandlerFunc {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	return func(c *gin.Context) {
		//!validation of the user's JWT token happens here
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Authentication Header is Required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Header is not in Proper format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid token shared"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Can not extract the claims"})
			c.Abort()
			return
		}

		c.Set("UserName", claims.Username)

		c.Next()
	}
}
