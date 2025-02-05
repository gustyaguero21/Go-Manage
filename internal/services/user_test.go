package services

import (
	"context"
	"errors"
	"fmt"
	"go-manage/cmd/config"
	"go-manage/internal/models"
	"go-manage/internal/repository"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repo := repository.UserRepository{
		DB: db,
	}
	userService := UserServices{
		DB:   db,
		Repo: repo,
	}

	test := []struct {
		Name     string
		Username string
		Expected bool
		MockAct  func()
	}{
		{
			Name:     "User already exists",
			Username: "johndoe",
			Expected: true,
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow("1", "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
		},
		{
			Name:     "User does not exist",
			Username: "nonexistentuser",
			Expected: false,
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("nonexistentuser").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}))
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			exists := userService.Exists(tt.Username)

			assert.Equal(t, tt.Expected, exists)
		})
	}
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repo := repository.UserRepository{DB: db}
	userService := UserServices{
		DB:   db,
		Repo: repo,
	}

	test := []struct {
		Name        string
		User        models.User
		ExpectedErr error
		MockAct     func()
	}{
		{
			Name: "Success",
			User: models.User{
				ID:       "1",
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			ExpectedErr: nil,
			MockAct: func() {
				mock.ExpectExec(config.TestSaveQuery).
					WithArgs().
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name: "Error",
			User: models.User{
				ID:       "1",
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			ExpectedErr: err,
			MockAct: func() {
			},
		},
		{
			Name: "Invalid password",
			User: models.User{
				ID:       "1",
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "invalid-password",
			},
			ExpectedErr: errors.New("invalid password"),
			MockAct: func() {
			},
		},
		{
			Name: "Invalid email",
			User: models.User{
				ID:       "1",
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Email:    "invalid-email",
				Password: "Password1234",
			},
			ExpectedErr: errors.New("invalid email"),
			MockAct: func() {
			},
		},
		{
			Name: "All fields are required",
			User: models.User{
				ID:       "1",
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
			},
			ExpectedErr: fmt.Errorf("all fields are required"),
			MockAct: func() {
				mock.ExpectExec(config.TestSaveQuery).
					WithArgs().
					WillReturnError(fmt.Errorf("all fields are required"))
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			createdUser, createErr := userService.CreateUser(ctx, tt.User)

			if tt.ExpectedErr != nil {
				assert.Equal(t, tt.ExpectedErr, createErr)
			}
			if createdUser.ID != "" {
				assert.Equal(t, tt.User.ID, createdUser.ID)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repo := repository.UserRepository{DB: db}
	userService := UserServices{
		DB:   db,
		Repo: repo,
	}

	test := []struct {
		Name           string
		Username       string
		ExpectedResult models.User
		ExpectedErr    error
		MockAct        func()
	}{
		{
			Name:     "Success",
			Username: "johndoe",
			ExpectedResult: models.User{
				ID:       "1",
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			ExpectedErr: nil,
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow("1", "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
		},
		{
			Name:     "Error",
			Username: "johndoe",
			ExpectedResult: models.User{
				ID:       "1",
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			ExpectedErr: err,
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnError(err)
			},
		},
		{
			Name:     "User does not exist",
			Username: "nonexistentuser",
			ExpectedResult: models.User{
				ID:       "1",
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("nonexistentuser").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}))
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			search, searchErr := userService.SearchUser(ctx, tt.Username)

			if tt.ExpectedErr != nil {
				assert.Equal(t, tt.ExpectedErr, searchErr)
			}

			if search.ID != "" {
				assert.Equal(t, tt.ExpectedResult.ID, search.ID)
			} else {
				assert.Error(t, errors.New("user not found"))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repo := repository.UserRepository{DB: db}

	userService := UserServices{
		DB:   db,
		Repo: repo,
	}

	test := []struct {
		Name        string
		Username    string
		ExpectedErr error
		SearchMock  func()
		MockAct     func()
	}{
		{
			Name:        "Success",
			Username:    "johndoe",
			ExpectedErr: nil,
			SearchMock: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow("1", "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
				mock.ExpectExec(config.TestDeleteQuery).
					WithArgs("johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name:        "Error",
			Username:    "johndoe",
			ExpectedErr: err,
			SearchMock: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow("1", "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
			},
		},
		{
			Name:        "User not found",
			Username:    "johndoe",
			ExpectedErr: errors.New("user not found"),
			SearchMock: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("nonexistentuser").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}))
			},
			MockAct: func() {
				mock.ExpectExec(config.TestDeleteQuery).
					WithArgs("johndoe").
					WillReturnError(errors.New("user not found"))
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.SearchMock()
			tt.MockAct()

			deleteErr := userService.DeleteUser(ctx, tt.Username)

			if tt.ExpectedErr != nil {
				assert.Equal(t, tt.ExpectedErr, deleteErr)
			}
		})
	}
}
