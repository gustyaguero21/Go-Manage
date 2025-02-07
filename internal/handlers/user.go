package handlers

import (
	"go-manage/internal/models"
	"go-manage/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gustyaguero21/go-core/pkg/web"
)

type UserHandler struct {
	userService services.UserServices
}

func NewUserHandler(userService services.UserServices) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (uh *UserHandler) Create(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		web.NewError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	created, createErr := uh.userService.CreateUser(ctx, user)
	if createErr != nil {
		web.NewError(ctx, http.StatusInternalServerError, createErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, created)
}

func (uh *UserHandler) Search(ctx *gin.Context) {
	username := ctx.Param("username")

	if username == "" {
		web.NewError(ctx, http.StatusBadRequest, "empty query param")
		return
	}

	search, searchErr := uh.userService.SearchUser(ctx, username)
	if searchErr != nil {
		web.NewError(ctx, http.StatusInternalServerError, searchErr.Error())
		return
	}

	if search.ID == "" {
		web.NewError(ctx, http.StatusOK, "user not found")
		return
	}

	ctx.JSON(http.StatusOK, search)

}
