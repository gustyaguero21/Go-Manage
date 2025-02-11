package handlers

import (
	"bytes"
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
			},
		},
		{
			Name:         "User not found",
			Username:     "johndoe",
			ExpectedCode: http.StatusNotFound,
			MockAct: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnError(config.ErrUserNotFound)
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

func TestCreate(t *testing.T) {
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
	r.POST("/create", handler.Create)

	tests := []struct {
		Name         string
		Body         string
		ExpectedCode int
		SearchMock   func()
		MockAct      func()
	}{
		{
			Name: "Success",
			Body: `{
				"id": "1",
				"name": "John",
				"surname": "Doe",
				"username": "johndoe",
				"email": "johndoe@example.com",
				"password": "Password1234"
			}`,
			ExpectedCode: http.StatusOK,
			SearchMock: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}))
			},
			MockAct: func() {
				mock.ExpectExec(config.TestSaveQuery).
					WithArgs().
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name:         "Invalid JSON",
			Body:         `{"id": "1", "name": "John", "surname": "Doe", "username": "johndoe", "email": "johndoe@example.com", "password": }`,
			ExpectedCode: http.StatusBadRequest,
			SearchMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name: "Error",
			Body: `{
				"id": "1",
				"name": "John",
				"surname": "Doe",
				"username": "johndoe",
				"email": "johndoe@example.com",
				"password": "Password1234"
			}`,
			ExpectedCode: http.StatusInternalServerError,
			SearchMock:   func() {},
			MockAct:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.SearchMock()
			tt.MockAct()

			req, _ := http.NewRequest(http.MethodPost, "/create", bytes.NewBufferString(tt.Body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
		})
	}
}

func TestDelete(t *testing.T) {
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
	r.DELETE("/delete", handler.Delete)

	tests := []struct {
		Name         string
		Username     string
		ExpectedCode int
		SearchMock   func()
		MockAct      func()
	}{
		{
			Name:         "Success",
			Username:     "johndoe",
			ExpectedCode: http.StatusOK,
			SearchMock: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
				mock.ExpectExec(config.TestDeleteQuery).
					WithArgs("johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name:         "Empty query param",
			Username:     "",
			ExpectedCode: http.StatusBadRequest,
			SearchMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name:         "Error",
			Username:     "johndoe",
			ExpectedCode: http.StatusInternalServerError,
			SearchMock: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(mock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow("1", "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
				mock.ExpectExec(config.TestDeleteQuery).
					WithArgs("johndoe").
					WillReturnError(err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.SearchMock()
			tt.MockAct()

			url := "/delete?username=" + tt.Username

			req, _ := http.NewRequest(http.MethodDelete, url, nil)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
		})
	}
}

func TestUpdate(t *testing.T) {
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
	r.PATCH("/update", handler.Update)

	tests := []struct {
		Name         string
		Username     string
		Body         string
		ExpectedCode int
		SearchMock   func()
		MockAct      func()
	}{
		{
			Name:     "Success",
			Username: "johndoe",
			Body: `{
				"id": "1",
				"name": "Johncito",
				"surname": "Doecito",
				"username": "johndoe",
				"email": "johndoe2024@example.com",
				"password": "NewPassword1234"
			}`,
			ExpectedCode: http.StatusOK,
			SearchMock: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
				mock.ExpectExec(config.TestUpdateQuery).
					WithArgs("Johncito", "Doecito", "johndoe2024@example.com", "johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name:         "Empty query param",
			Username:     "",
			Body:         `{"name": "Johncito", "surname": "Doecito", "email": "johndoe2024@example.com", "password": "NewPassword1234"}`,
			ExpectedCode: http.StatusBadRequest,
			SearchMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name:         "Invalid JSON",
			Username:     "johndoe",
			Body:         `{"name": "Johncito", "surname": "Doecito", "email": "johndoe2024@example.com", "password":}`,
			ExpectedCode: http.StatusBadRequest,
			SearchMock:   func() {},
			MockAct:      func() {},
		},
		{
			Name:         "Error",
			Username:     "johndoe",
			Body:         `{"name": "Johncito", "surname": "Doecito", "email": "johndoe2024@example.com", "password": "NewPassword1234"}`,
			ExpectedCode: http.StatusInternalServerError,
			SearchMock: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
				mock.ExpectExec(config.TestUpdateQuery).
					WithArgs("Johncito", "Doecito", "johndoe2024@example.com", "johndoe").
					WillReturnError(errors.New("update error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.SearchMock()
			tt.MockAct()

			body := bytes.NewBufferString(tt.Body)
			url := "/update?username=" + tt.Username
			req, _ := http.NewRequest(http.MethodPatch, url, body)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
		})
	}
}

func TestChangePwd(t *testing.T) {
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
	r.PATCH("/change-password", handler.ChangePwd)

	test := []struct {
		Name         string
		Username     string
		NewPassword  string
		ExpectedCode int
		SearchMock   func()
		MockAct      func()
	}{
		{
			Name:         "Success",
			Username:     "johndoe",
			NewPassword:  "NewPassword1234",
			ExpectedCode: http.StatusOK,
			SearchMock: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
				mock.ExpectExec(config.TestChangePwdQuery).
					WithArgs(sqlmock.AnyArg(), "johndoe").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name:         "Empty query param",
			Username:     "",
			NewPassword:  "",
			ExpectedCode: http.StatusBadRequest,
			SearchMock: func() {
			},
			MockAct: func() {},
		},
		{
			Name:         "Error",
			Username:     "johndoe",
			NewPassword:  "NewPassword1234",
			ExpectedCode: http.StatusInternalServerError,
			SearchMock: func() {
				mock.ExpectQuery(config.TestSearchQuery).
					WithArgs("johndoe").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "username", "email", "password"}).
						AddRow(1, "John", "Doe", "johndoe", "johndoe@example.com", "Password1234"))
			},
			MockAct: func() {
				mock.ExpectExec(config.TestChangePwdQuery).
					WithArgs(sqlmock.AnyArg(), "johndoe").
					WillReturnError(err)
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.SearchMock()
			tt.MockAct()

			url := "/change-password?username=" + tt.Username + "&new_password=" + tt.NewPassword

			req, _ := http.NewRequest(http.MethodPatch, url, nil)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
		})
	}
}
