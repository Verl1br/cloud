package repository

import (
	"fmt"

	"github.com/dhevve/uploadImage/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int

	createQueryUser := fmt.Sprintf("INSERT INTO %s (username, password_hash) values ($1, $2) RETURNING id", "users")

	row := r.db.QueryRow(createQueryUser, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", "users")
	err := r.db.Get(&user, query, username, password)
	return user, err
}
