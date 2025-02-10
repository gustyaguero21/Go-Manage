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
		web.NewError(ctx, http.StatusBadRequest, config.ErrEmptyQueryParam)
		return
	}

	search, searchErr := h.userService.SearchUser(ctx, username)
	if searchErr != nil {
		if strings.Contains(searchErr.Error(), config.ErrUserNotFound) {
			web.NewError(ctx, http.StatusNotFound, config.ErrUserNotFound)
			return
		} else {
			web.NewError(ctx, http.StatusInternalServerError, searchErr.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, searchResponse(config.SuccessStatus, config.SearchMessage, search))

}

func searchResponse(status string, message string, user models.User) *models.SearchResponse {
	return &models.SearchResponse{
		Status:  status,
		Message: message,
		User:    user,
	}
}
