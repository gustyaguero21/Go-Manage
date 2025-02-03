package repository

import (
	"database/sql"
	"fmt"
	"go-manage/cmd/config"
	"go-manage/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (ur *UserRepository) Exists(existsQuery, username string) bool {

	search, searchErr := ur.Search(config.SearchUserQuery, username)
	if searchErr != nil {
		return false
	}

	if search.ID != "" {
		return true
	}
	return false
}

func (ur *UserRepository) Search(searchQuery, username string) (models.User, error) {
	user := models.User{}
	rows, err := ur.DB.Query(config.SearchUserQuery, username)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return models.User{}, fmt.Errorf("user not found")
		}
	}

	return user, nil
}

func (ur *UserRepository) Save(saveQuery string, user models.User) error {
	_, saveErr := ur.DB.Exec(saveQuery, user.ID, user.Name, user.Surname, user.Username, user.Email, user.Password)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func (ur *UserRepository) Delete(deleteQuery, username string) error {
	_, deleteErr := ur.DB.Exec(deleteQuery, username)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}
