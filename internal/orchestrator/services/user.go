package services

import (
	"errors"
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/repositories"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	repo *repositories.User
}

func (u *User) Create(login, password string) (int, error) {
	_, err := u.repo.GetByLogin(login)
	if err == nil {
		// user exists with current login
		return 0, errors.New("user with login is exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	userId, err := u.repo.Create(login, string(hash))

	return userId, err
}

func (u *User) Authorization(login, password string) int {
	user, err := u.repo.GetByLogin(login)
	if err != nil {
		return 0
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0
	}
	return user.UserID
}

// create new user service
func UserService() *User { return &User{} }
