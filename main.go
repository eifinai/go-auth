package main

import (
	"fmt"
	"go-auth/controllers"
	"go-auth/middleware"
	"log"
	"os"
	"strconv"

	"go-auth/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	postgresPassword = "POSTGRES_PASSWORD"
	host             = "HOST"
	port             = "PORT"
	user             = "DB_USER"
	dbname           = "DB_NAME"
)

type options struct {
	postgresPassword string
	host             string
	port             int
	user             string
	dbname           string
}

func getOptions() (options, error) {
	port, err := strconv.Atoi(os.Getenv(port))

	if err != nil {
		return options{}, err
	}

	return options{
		postgresPassword: os.Getenv(postgresPassword),
		host:             os.Getenv(host),
		port:             port,
		user:             os.Getenv(user),
		dbname:           os.Getenv(dbname),
	}, nil
}
func getPostgresString(opt options) string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		opt.host, opt.port, opt.user, opt.postgresPassword, opt.dbname)
}

func init() {

}

func main() {
	opt, err := getOptions()
	if err != nil {
		log.Fatal(err)
	}
	psqlInfo := getPostgresString(opt)
	dbConnection, err := database.NewDatabase(psqlInfo)

	if err != nil {
		panic(err)
	}
	router := controllers.Routes{DB: dbConnection}
	middleware := middleware.MiddlewareAuth{R: router}

	fmt.Println(dbConnection)
	fmt.Println("Successfully connected!")
	fmt.Println(psqlInfo)

	r := gin.Default()

	r.POST("/signup", router.Signup)
	r.POST("/login", router.Login)
	r.GET("/validate", middleware.RequireAuth, router.Validate)

	r.Run(":8080")
}
