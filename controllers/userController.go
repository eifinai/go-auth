package controllers

import (
	"fmt"
	"net/http"

	"go-auth/database"
	"go-auth/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Routes struct {
	DB database.Database
}

func (r Routes) Signup(c *gin.Context) {
	// Get Email/Password from request body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Check if user already exists

	// Create user in the database
	user := models.User{Email: body.Email, Password: string(hash)}
	_, err = r.DB.Exec(fmt.Sprintf("INSERT INTO users (username, password_hash) VALUES ('%s', '%s')", user.Email, user.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": string(hash)})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{"message": "Successfully signed up"})
}
