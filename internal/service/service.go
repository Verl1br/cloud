package service

import (
	"github.com/dhevve/uploadImage/internal/models"
	"github.com/dhevve/uploadImage/internal/repository"
	"github.com/gin-gonic/gin"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	ParseToken(accessToken string) (int, error)
	GenerateToken(username, password string) (string, error)
}

type UploadImage interface {
	Upload(c *gin.Context, userId int) (int, error)
	GetAll(userId int) ([]models.Image, error)
	GetById(userId, imageId int) (models.Image, error)
	Delete(userId, itemId int) error
}

type Service struct {
	Authorization
	UploadImage
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
		UploadImage:   NewUploadService(repo),
	}
}
