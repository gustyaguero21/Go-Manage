package handlers

import (
	"errors"
	"go-manage/cmd/config"
	"go-manage/internal/repository"
	"go-manage/internal/services"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestSearch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.UserRepository{DB: db}
	userService := services.UserServices{DB: db, Repo: repo}
	handler := UserHandler{userService: userService}

	r := gin.Default()
	r.GET("/search", handler.Search)

	tests := []struct {
		Name         string
		Username     string
		ExpectedCode int
		MockAct      func()
	}{
		{
			Name:         "Success",
			Username:     "johndoe",
			ExpectedCode: http.StatusOK,
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
		},
		{
			Name:         "Error",
			Username:     "johndoe",
			ExpectedCode: http.StatusInternalServerError,
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnError(err)
			},
		},
		{
			Name:         "Empty query param",
			Username:     "",
			ExpectedCode: http.StatusBadRequest,
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("").
					WillReturnError(errors.New(config.ErrEmptyQueryParam))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockAct()

			url := "/search?username=" + tt.Username
			req, _ := http.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)

		})
	}
}
