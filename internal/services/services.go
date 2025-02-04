package services

import (
	"context"
	"go-manage/internal/models"
)

type Services interface {
	CreateUser(ctx context.Context, user models.User) (created models.User, err error)
}
