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

func (us *UserServices) SearchUser(ctx context.Context, username string) (search models.User, err error) {
	search, searchErr := us.Repo.Search(config.SearchUserQuery, username)
	if searchErr != nil {
		return models.User{}, errors.New("error searching user. Error: " + searchErr.Error())
	}

	if search.ID == "" {
		return models.User{}, errors.New("user not found")
	}

	return search, nil
}

func (us *UserServices) DeleteUser(ctx context.Context, username string) (err error) {
	if !us.Exists(username) {
		return errors.New("user not found")
	}

	if deleteErr := us.Repo.Delete(config.DeleteUserQuery, username); deleteErr != nil {
		return errors.New("error deleting user. Error: " + deleteErr.Error())
	}
	return nil
}

func (us *UserServices) UpdateUser(ctx context.Context, username string, user models.User) (updated models.User, err error) {
	if !us.Exists(username) {
		return models.User{}, errors.New("user not found")
	}

	if updateErr := us.Repo.Update(config.UpdateUserQuery, username, user); updateErr != nil {
		return models.User{}, errors.New("error updating user. Error: " + updateErr.Error())
	}

	return user, nil
}

func (us *UserServices) ChangeUserPwd(ctx context.Context, username string, newPassword string) (err error) {
	if !us.Exists(username) {
		return errors.New("user not found")
	}

	if changePwd := us.Repo.ChangePwd(config.ChangeUserPwdQuery, username, newPassword); changePwd != nil {
		return errors.New("error changing user password. Error: " + changePwd.Error())
	}

	return nil
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
