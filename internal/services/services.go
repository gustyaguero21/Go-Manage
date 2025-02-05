package services

import (
	"context"
	"go-manage/internal/models"
)

type Services interface {
	Exists(username string) error
	CreateUser(ctx context.Context, user models.User) (created models.User, err error)
	SearchUser(ctx context.Context, username string) (search models.User, err error)
	DeleteUser(ctx context.Context, username string) (err error)
	UpdateUser(ctx context.Context, username string, user models.User) (updated models.User, err error)
}
