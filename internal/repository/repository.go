package repository

import (
	"github.com/dhevve/uploadImage/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type UploadImage interface {
	Upload(userId int, fullFileName string) (int, error)
	GetAll(userId int) ([]models.Image, error)
	GetById(userId, imageId int) (models.Image, error)
	Delete(userId, itemId int) error
}

type Repository struct {
	Authorization
	UploadImage
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UploadImage:   NewUploadPostgres(db),
		Authorization: NewAuthPostgres(db),
	}
}
