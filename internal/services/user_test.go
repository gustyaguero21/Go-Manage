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

func TestCreate(t *testing.T) {
	ctx := context.Background()
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
			Expected: true, // Aquí esperamos que Exists devuelva `true` porque el usuario existe
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
			Expected: false, // Aquí esperamos que Exists devuelva `false` porque el usuario no existe
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

			// Verificamos si el usuario existe usando el servicio
			exists := userService.Exists(tt.Username)

			// Comprobamos que el valor de exists sea el esperado
			assert.Equal(t, tt.Expected, exists)
		})
	}
}
