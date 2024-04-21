package repositories

import (
	"context"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/database"
	"time"
)

type User struct {
}

type UserModel struct {
	UserID    int       `json:"user_id"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Find user by id
func (u *User) GetById(userId int) (UserModel, error) {
	var user UserModel

	query := "SELECT * FROM users WHERE user_id = $1"
	if err := database.DB.QueryRow(context.Background(), query, userId).Scan(&user.UserID, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return user, err
	}

	return user, nil
}

// Find user by login
func (u *User) GetByLogin(login string) (UserModel, error) {
	var user UserModel

	query := "SELECT * FROM users WHERE login = $1"
	if err := database.DB.QueryRow(context.Background(), query, login).Scan(&user.UserID, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return user, err
	}

	return user, nil
}

// Create new row with user
func (u *User) Create(login string, password string) (int, error) {
	var insertedID int

	query := "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING user_id"
	if err := database.DB.QueryRow(context.Background(), query, login, password).Scan(&insertedID); err != nil {
		return 0, err
	}

	return insertedID, nil
}

func UserRepository() *User {
	return &User{}
}
