package controllers

import (
	"net/http"
	"os"

	"go-auth/database"
	"go-auth/models"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	var existingUser models.User
	r.DB.DB.QueryRow("SELECT * FROM users WHERE username = $1", body.Email).Scan(&existingUser.Id, &existingUser.Email, &existingUser.Password)

	if existingUser.Id != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This user already exists"})
		return
	}

	// Create user in the database
	user := models.User{Email: body.Email, Password: string(hash)}
	_, err = r.DB.DB.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", user.Email, user.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{"message": "Successfully signed up"})
}
func (r Routes) Login(c *gin.Context) {
	//Get email and password from request body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	//Get user from database
	var user models.User
	r.DB.DB.QueryRow("SELECT * FROM users WHERE username = $1", body.Email).Scan(&user.Id, &user.Email, &user.Password)

	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password or email"})
		return
	}
	//Compare password with hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password or email"})
		return
	}

	//Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(os.Getenv("SECRET"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	//respond
	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}
