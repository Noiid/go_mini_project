package main

import (
	"database/sql"
	"jwt/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) FetchUser() ([]model.User, error) {
	preparedStatement := `
	SELECT 
	    *
	FROM 
		users`

	rows, err := r.db.Query(preparedStatement)
	if err != nil {
		return nil, err
	}

	var users []model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
