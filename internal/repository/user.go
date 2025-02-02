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
	if user.ID == "" {
		return models.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}
