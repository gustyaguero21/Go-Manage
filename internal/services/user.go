package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-manage/cmd/config"
	"go-manage/internal/models"
	"go-manage/internal/repository"

	"github.com/google/uuid"
	"github.com/gustyaguero21/go-core/pkg/encrypter"
	"github.com/gustyaguero21/go-core/pkg/validator"
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

	user.ID = uuid.New().String()

	hashedPwd, hashErr := encrypter.PasswordEncrypter(user.Password)
	if hashErr != nil {
		return models.User{}, hashErr
	}

	user.Password = string(hashedPwd)

	if createErr := us.Repo.Save(config.SaveUserQuery, user); createErr != nil {
		return models.User{}, errors.New(config.ErrSaveUser + createErr.Error())
	}

	return user, nil
}

func (us *UserServices) SearchUser(ctx context.Context, username string) (search models.User, err error) {
	search, searchErr := us.Repo.Search(config.SearchUserQuery, username)
	if searchErr != nil {
		return models.User{}, errors.New(config.ErrSearchUser + searchErr.Error())
	}

	return search, nil
}

func (us *UserServices) DeleteUser(ctx context.Context, username string) (err error) {
	if !us.Exists(username) {
		return errors.New(config.ErrUserNotFound)
	}

	if deleteErr := us.Repo.Delete(config.DeleteUserQuery, username); deleteErr != nil {
		return errors.New(config.ErrDeleteUser + deleteErr.Error())
	}
	return nil
}

func (us *UserServices) UpdateUser(ctx context.Context, username string, user models.User) (update models.User, err error) {
	if !us.Exists(username) {
		return models.User{}, errors.New(config.ErrUserNotFound)
	}

	search, searchErr := us.Repo.Search(config.SearchUserQuery, username)
	if searchErr != nil {
		return models.User{}, searchErr
	}

	if updateErr := us.Repo.Update(config.UpdateUserQuery, username, user); updateErr != nil {
		return models.User{}, errors.New(config.ErrUpdateUser + updateErr.Error())
	}

	updated := models.User{
		ID:       search.ID,
		Name:     user.Name,
		Surname:  user.Surname,
		Username: search.Username,
		Email:    search.Email,
		Password: search.Password,
	}

	return updated, nil
}

func (us *UserServices) ChangeUserPwd(ctx context.Context, username string, newPassword string) (err error) {
	if !us.Exists(username) {
		return errors.New(config.ErrUserNotFound)
	}

	if changePwd := us.Repo.ChangePwd(config.ChangeUserPwdQuery, username, newPassword); changePwd != nil {
		return errors.New(config.ErrChangePwd + changePwd.Error())
	}

	return nil
}

func paramsValidation(user models.User) error {
	if user.Name == "" || user.Surname == "" || user.Username == "" || user.Email == "" || user.Password == "" {
		return fmt.Errorf(config.ErrAllFieldsAreRequired)
	}
	if !validator.ValidatePassword(user.Password) {
		return fmt.Errorf(config.ErrInvalidPassword)
	}

	if !validator.ValidateEmail(user.Email) {
		return fmt.Errorf(config.ErrInvalidPassword)
	}

	return nil
}
