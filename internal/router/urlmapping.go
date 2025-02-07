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
		log.Fatal(connErr)
	}

	repo := repository.UserRepository{DB: conn}
	userService := services.UserServices{
		DB:   conn,
		Repo: repo,
	}

	handler := handlers.NewUserHandler(userService)

	api := r.Group("/api/go-manage")

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	api.POST("/create", handler.Create)
	api.GET("/search/:username", handler.Search)
	api.DELETE("/delete/:username", handler.Delete)
	api.PATCH("/update/:username", handler.Update)
	api.PATCH("/change-password/:username", handler.ChangePwd)

}
