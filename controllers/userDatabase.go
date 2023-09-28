package controllers

import (
	"database/sql"
	"go-auth/models"
)

type UserDatabase struct {
	DB *sql.DB
}

func (u UserDatabase) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := u.DB.QueryRow("SELECT * FROM users WHERE username = $1", email).Scan(&user.Id, &user.Email, &user.Password)

	return user, err
}
func (u UserDatabase) CreateUser(user models.User) error {
	_, err := u.DB.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", user.Email, user.Password)

	return err
}
func (u UserDatabase) GetUserById(id int) (models.User, error) {
	var user models.User
	err := u.DB.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.Id, &user.Email, &user.Password)

	return user, err
}
