package handlers

import (
	"bytes"
	"encoding/json"
	"go-manage/cmd/config"
	"go-manage/internal/models"
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

	test := []struct {
		Name         string
		User         models.User
		ExpectedCode int
		SearchMock   func()
		MockAct      func()
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
			Name: "Error",
			User: models.User{
				ID:       "1",
				Name:     "John",
				Surname:  "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			ExpectedCode: http.StatusInternalServerError,
			SearchMock: func() {
			},
			MockAct: func() {
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.SearchMock()
			tt.MockAct()

			body, err := json.Marshal(tt.User)
			if err != nil {
				t.Fatal("Error marshaling user:", err)
			}

			req, _ := http.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(body))

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

	test := []struct {
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
			SearchMock: func() {
			},
			MockAct: func() {},
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

	for _, tt := range test {
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

	test := []struct {
		Name         string
		Username     string
		Update       models.User
		ExpectedCode int
		SearchMock   func()
		MockAct      func()
	}{
		{
			Name:     "Success",
			Username: "johndoe",
			Update: models.User{
				ID:       "1",
				Name:     "Johncito",
				Surname:  "Doecito",
				Username: "johndoe",
				Email:    "johndoe2024@example.com",
				Password: "NewPassword1234",
			},
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
			Name:     "Error",
			Username: "johndoe",
			Update: models.User{
				ID:       "1",
				Name:     "Johncito",
				Surname:  "Doecito",
				Username: "johndoe",
				Email:    "johndoe2024@example.com",
				Password: "NewPassword1234",
			},
			ExpectedCode: http.StatusInternalServerError,
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
			Name:         "Empty query param",
			Username:     "",
			Update:       models.User{},
			ExpectedCode: http.StatusBadRequest,
			SearchMock: func() {
			},
			MockAct: func() {},
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			tt.SearchMock()
			tt.MockAct()

			body, err := json.Marshal(tt.Update)
			if err != nil {
				t.Fatal("Error marshaling user:", err)
			}

			url := "/update?username=" + tt.Username

			req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))

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
