package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-manage/cmd/config"
	"go-manage/internal/models"
	"go-manage/internal/repository"

	"github.com/gustyaguero21/Go-Core/pkg/encrypter"
	"github.com/gustyaguero21/Go-Core/pkg/validator"
)

type UserServices struct {
	DB   *sql.DB
	Repo repository.UserRepository
}

func (us *UserServices) Exists(username string) bool {
	exists := us.Repo.Exists(config.SearchUserQuery, username)

	return exists
}

func (us *UserServices) CreateUser(ctx context.Context, user models.User) (created models.User, err error) {
	if us.Exists(user.Username) {
		return models.User{}, errors.New("user already exists")
	}
	if checkErr := paramsValidation(user); checkErr != nil {
		return models.User{}, checkErr
	}

	hashedPwd, hashErr := encrypter.PasswordEncrypter(user.Password)
	if hashErr != nil {
		return models.User{}, hashErr
	}

	user.Password = string(hashedPwd)

	if createErr := us.Repo.Save(config.SaveUserQuery, user); createErr != nil {
		return models.User{}, errors.New("error creating user. Error: " + createErr.Error())
	}

	return user, nil
}

func paramsValidation(user models.User) error {
	if user.ID == "" || user.Name == "" || user.Surname == "" || user.Username == "" || user.Email == "" || user.Password == "" {
		return fmt.Errorf("all fields are required")
	}
	if !validator.ValidatePassword(user.Password) {
		return fmt.Errorf("invalid password")
	}

	if !validator.ValidateEmail(user.Email) {
		return fmt.Errorf("invalid email")
	}

	return nil
}
