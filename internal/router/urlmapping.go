package router

import (
	"go-manage/internal/data"
	"go-manage/internal/handlers"
	"go-manage/internal/repository"
	"go-manage/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Urlmapping(r *gin.Engine) {

	conn, connErr := data.InitDatabase()
	if connErr != nil {
		log.Fatal("cannot initialize database. Error: " + connErr.Error())
	}

	repo := repository.UserRepository{DB: conn}
	userService := services.UserServices{DB: conn, Repo: repo}

	handler := handlers.NewUserHandler(userService)

	api := r.Group("/api/go-manage")

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	api.GET("/search", handler.Search)
	api.POST("/create", handler.Create)
	api.DELETE("/delete", handler.Delete)
	api.PATCH("/update", handler.Update)
	api.PATCH("/change-password", handler.ChangePwd)
}
