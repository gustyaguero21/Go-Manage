package repository

import (
	"fmt"
	"go-manage/cmd/config"
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

	repo := UserRepository{
		DB: db,
	}

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
				mock.ExpectQuery(config.SearchQuery).
					WithArgs("johndoe2024").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe2024", "john@example.com", "password123"))
			},
		},
		{
			Name:          "User not found",
			Username:      "johndoe2024",
			ExpectedError: fmt.Errorf("user not found"),
			MockAct: func() {
				mock.ExpectQuery(config.SearchQuery).
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
				mock.ExpectQuery(config.SearchQuery).
					WithArgs("johndoe2024").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}))
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			search, searchErr := repo.Search(config.SearchQuery, tt.Username)

			if tt.ExpectedError != nil {
				assert.Equal(t, tt.ExpectedError, searchErr)
			}
			if search.ID != "" {
				assert.Error(t, fmt.Errorf("user not found"))
			}
		})
	}
}
