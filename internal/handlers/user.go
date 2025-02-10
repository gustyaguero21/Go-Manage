package handlers

import (
	"go-manage/cmd/config"
	"go-manage/internal/models"
	"go-manage/internal/services"
	"net/http"
	"strings"

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

func (h *UserHandler) Search(ctx *gin.Context) {
	username := ctx.Query("username")

	if username == "" {
		web.NewError(ctx, http.StatusBadRequest, config.ErrEmptyQueryParam.Error())
		return
	}

	search, searchErr := h.userService.SearchUser(ctx, username)
	if searchErr != nil {
		if strings.Contains(searchErr.Error(), config.ErrUserNotFound.Error()) {
			web.NewError(ctx, http.StatusNotFound, config.ErrUserNotFound.Error())
			return
		} else {
			web.NewError(ctx, http.StatusInternalServerError, searchErr.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, searchResponse(config.SuccessStatus, config.SearchMessage, search))

}

func (h *UserHandler) Create(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		web.NewError(ctx, http.StatusBadRequest, err.Error())
		return
	}

}

func searchResponse(status string, message string, user models.User) *models.SearchResponse {
	return &models.SearchResponse{
		Status:  status,
		Message: message,
		User:    user,
	}
}
