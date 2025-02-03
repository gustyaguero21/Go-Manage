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

func TestExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := UserRepository{DB: db}

	tests := []struct {
		Name         string
		Username     string
		ExpectedBool bool
		MockAct      func()
	}{
		{
			Name:         "Success",
			Username:     "johndoe",
			ExpectedBool: true,
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "password123"))
			},
		},
		{
			Name:         "Error",
			Username:     "johndoe2024",
			ExpectedBool: false,
			MockAct: func() {
			},
		},
		{
			Name:         "User not found",
			Username:     "lala",
			ExpectedBool: false,
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("lala").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			exists := repo.Exists(config.TestSearchQuery, tt.Username)

			assert.Equal(t, tt.ExpectedBool, exists)
		})
	}
}

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

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
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
			Username:      "johndoe",
			ExpectedError: nil,
			MockAct: func() {
				mock.ExpectExec(config.TestDeleteQuery).
					WithArgs("johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name:          "Error",
			Username:      "johndoe",
			ExpectedError: err,
			MockAct: func() {
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			deleteErr := repo.Delete(config.TestDeleteQuery, tt.Username)

			if tt.ExpectedError != nil {
				assert.Equal(t, tt.ExpectedError, deleteErr)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := UserRepository{DB: db}

	test := []struct {
		Name        string
		Username    string
		User        models.User
		ExpectedErr error
		MockAct     func()
	}{
		{
			Name:     "Success",
			Username: "johndoe",
			User: models.User{
				ID:       "1",
				Name:     "Johncito",
				Surname:  "Doecito",
				Username: "johndoe",
				Email:    "johndoe2024@example.com",
				Password: "Password1234",
			},
			ExpectedErr: nil,
			MockAct: func() {
				mock.ExpectExec(config.TestUpdateQuery).
					WithArgs("Johncito", "Doecito", "johndoe2024@example.com", "johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name:     "Error",
			Username: "johndoe",
			User: models.User{
				ID:       "1",
				Name:     "Johncito",
				Surname:  "Doecito",
				Username: "johndoe",
				Email:    "johndoe2024@example.com",
				Password: "Password1234",
			},
			ExpectedErr: fmt.Errorf("error updating user"),
			MockAct: func() {
				mock.ExpectExec(config.TestUpdateQuery).
					WithArgs("Johncito", "Doecito", "johndoe2024@example.com", "johndoe").
					WillReturnError(fmt.Errorf("error updating user"))
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			updateErr := repo.Update(config.UpdateUserQuery, tt.Username, tt.User)

			if tt.ExpectedErr != nil {
				assert.Equal(t, tt.ExpectedErr.Error(), updateErr.Error())
			} else {
				assert.NoError(t, updateErr)
			}
		})
	}
}
