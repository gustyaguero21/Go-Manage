package services

import (
	"context"
	"go-manage/internal/models"
)

type Services interface {
	Exists(username string) error
	CreateUser(ctx context.Context, user models.User) (created models.User, err error)
}
