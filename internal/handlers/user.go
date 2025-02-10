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

	created, createErr := h.userService.CreateUser(ctx, user)
	if createErr != nil {
		web.NewError(ctx, http.StatusInternalServerError, createErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, createResponse(config.SuccessStatus, config.CreateMessage, created))
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	username := ctx.Query("username")
	if username == "" {
		web.NewError(ctx, http.StatusBadRequest, config.ErrEmptyQueryParam.Error())
		return
	}

	if deleteErr := h.userService.DeleteUser(ctx, username); deleteErr != nil {
		web.NewError(ctx, http.StatusInternalServerError, deleteErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, deleteResponse(config.SuccessStatus, config.DeleteMessage))
}

func (h *UserHandler) Update(ctx *gin.Context) {
	username := ctx.Query("username")
	if username == "" {
		web.NewError(ctx, http.StatusBadRequest, config.ErrEmptyQueryParam.Error())
		return
	}

	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		web.NewError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if updateErr := h.userService.UpdateUser(ctx, username, user); updateErr != nil {
		web.NewError(ctx, http.StatusInternalServerError, updateErr.Error())
		return
	}
	ctx.JSON(http.StatusOK, updateResponse(config.SuccessStatus, config.UpdateMessage))
}

func (h *UserHandler) ChangePwd(ctx *gin.Context) {
	username := ctx.Query("username")
	newPassword := ctx.Query("new_password")
	if username == "" || newPassword == "" {
		web.NewError(ctx, http.StatusBadRequest, config.ErrEmptyQueryParam.Error())
		return
	}

	if changeErr := h.userService.ChangeUserPwd(ctx, username, newPassword); changeErr != nil {
		web.NewError(ctx, http.StatusInternalServerError, changeErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, changePwdResponse(config.SuccessStatus, config.ChangePwdMessage))
}

func searchResponse(status string, message string, user models.User) *models.SearchResponse {
	return &models.SearchResponse{
		Status:  status,
		Message: message,
		User:    user,
	}
}

func createResponse(status string, message string, created models.User) *models.CreateUserResponse {
	return &models.CreateUserResponse{
		Status:  status,
		Message: message,
		Created: created,
	}
}

func deleteResponse(status string, message string) *models.DeleteUserResponse {
	return &models.DeleteUserResponse{
		Status:  status,
		Message: message,
	}
}

func updateResponse(status string, message string) *models.UpdateUserResponse {
	return &models.UpdateUserResponse{
		Status:  status,
		Message: message,
	}
}

func changePwdResponse(status string, message string) *models.ChangePwdResponse {
	return &models.ChangePwdResponse{
		Status:  status,
		Message: message,
	}
}
