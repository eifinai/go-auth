package controllers

import (
	"database/sql"
	"net/http"
	"os"

	"go-auth/models"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Routes struct {
	DB     *sql.DB
	UserDB UserDatabase
}

func NewRouter(db *sql.DB) Routes {
	return Routes{
		DB:     db,
		UserDB: UserDatabase{DB: db},
	}
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
	existingUser, _ = r.UserDB.GetUserByEmail(body.Email)
	//TODO: check for no result error
	if existingUser.Id != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	// Create user in the database
	user := models.User{Email: body.Email, Password: string(hash)}
	err = r.UserDB.CreateUser(user)
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
	user, _ = r.UserDB.GetUserByEmail(body.Email)
	//TODO: check for no result error
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

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	//respond
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("accessToken", tokenString, 3600*24, "", "localhost", false, true)
	//c.JSON(http.StatusOK, gin.H{"accessToken": tokenString})

}
func (r Routes) Validate(c *gin.Context) {
	//user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully validated"})
}
