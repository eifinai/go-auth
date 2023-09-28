package middleware

import (
	"go-auth/controllers"
	"go-auth/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MiddlewareAuth struct {
	Routes controllers.Routes
}

func (m MiddlewareAuth) RequireAuth(c *gin.Context) {
	//Get token from request header
	tokenString, err := c.Cookie("accessToken")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//Validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//Check if token is expired
		if claims["exp"].(float64) < float64(time.Now().Unix()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		}

		//find user in database
		var user models.User
		//m.Routes.DB.QueryRow("SELECT * FROM users WHERE id = $1", claims["sub"]).Scan(&user.Id, &user.Email, &user.Password)
		user, err = m.Routes.UserDB.GetUserById(int(claims["sub"].(float64)))

		if err != nil || user.Id == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//add user to context
		c.Set("user", user)

		//call next middleware
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
