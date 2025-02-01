package router

import (
	"go-manage/internal/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Urlmapping(r *gin.Engine) {
	api := r.Group("/api/go-manage")
	data.InitDatabase()

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})
}
