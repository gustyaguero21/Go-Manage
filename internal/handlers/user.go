package handlers

import (
	"go-manage/internal/services"
)

type UserHandler struct {
	userService services.UserServices
}

func NewUserHandler(userService services.UserServices) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
