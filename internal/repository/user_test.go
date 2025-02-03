package repository

import (
	"fmt"
	"go-manage/cmd/config"
	"go-manage/internal/models"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := UserRepository{DB: db}

	test := []struct {
		Name          string
		Username      string
		ExpectedError error
		MockAct       func()
	}{
		{
			Name:          "Success",
			Username:      "johndoe2024",
			ExpectedError: nil,
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe2024").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe2024", "john@example.com", "password123"))
			},
		},
		{
			Name:          "Error",
			Username:      "johndoe2024",
			ExpectedError: err,
			MockAct: func() {
			},
		},
		{
			Name:          "User not found",
			Username:      "johndoe2024",
			ExpectedError: fmt.Errorf("user not found"),
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe2024").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(nil, nil, nil, nil, nil, nil))
			},
		},
		{
			Name:          "Empty ID",
			Username:      "johndoe2024",
			ExpectedError: fmt.Errorf("user not found"),
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe2024").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}))
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			search, searchErr := repo.Search(config.TestSearchQuery, tt.Username)

			if tt.ExpectedError != nil {
				assert.Equal(t, tt.ExpectedError, searchErr)
			}
			if search.ID != "" {
				assert.Error(t, fmt.Errorf("user not found"))
			}
		})
	}
}

func TestSave(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := UserRepository{DB: db}

	test := []struct {
		Name          string
		User          models.User
		ExpectedError error
		MockAct       func()
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
			ExpectedError: nil,
			MockAct: func() {
				mock.ExpectExec(config.TestSaveQuery).
					WithArgs("1", "John", "Doe", "johndoe", "johndoe@example.com", "Password1234").
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
			ExpectedError: err,
			MockAct: func() {
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			saveErr := repo.Save(config.TestSaveQuery, tt.User)

			if tt.ExpectedError != nil {
				assert.Equal(t, tt.ExpectedError, saveErr)
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

	repo := UserRepository{DB: db}

	test := []struct {
		Name             string
		Username         string
		ExpectedResponse bool
		MockAct          func()
	}{
		{
			Name:             "Error",
			Username:         "johndoe",
			ExpectedResponse: false,
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe2024").
					WillReturnRows(mock.NewRows([]string{"username"}))
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			exists := repo.Exists(config.TestExistsQuery, tt.Username)

			assert.Equal(t, tt.ExpectedResponse, exists)
		})
	}
}
