package repository

import "go-manage/internal/models"

type Repository interface {
	Exists(existsQuery, username string) bool
	Search(searchQuery, username string) (models.User, error)
	Save(saveQuery string, user models.User) (models.User, error)
	Delete(deleteQuery, username string) error
	Update(updateQuery, username string, user models.User) error
	ChangePwd(changePwdQuery, username, newPassword string) error
}
