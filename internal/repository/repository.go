package repository

import "go-manage/internal/models"

type Repository interface {
	Search(searchQuery, username string) (models.User, error)
	Save(saveQuery string, user models.User) (models.User, error)
}
