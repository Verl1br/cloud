package service

import (
	"fmt"
	"os"

	"github.com/dhevve/uploadImage/internal/models"
	"github.com/dhevve/uploadImage/internal/repository"
	"github.com/gin-gonic/gin"
)

type UploadService struct {
	repo repository.UploadImage
}

func NewUploadService(repo repository.UploadImage) *UploadService {
	return &UploadService{repo: repo}
}

func (s *UploadService) Upload(c *gin.Context, userId int) (int, error) {
	fullFileName, err := saveFile(c)
	if err != nil {
		return 0, err
	}
	return s.repo.Upload(userId, fullFileName)
}

func (s *UploadService) GetAll(userId int) ([]models.Image, error) {
	return s.repo.GetAll(userId)
}

func (s *UploadService) GetById(userId, imageId int) (models.Image, error) {
	return s.repo.GetById(userId, imageId)
}

func (s *UploadService) Delete(userId, itemId int) error {
	image, err := s.GetById(userId, itemId)
	if err != nil {
		return err
	}
	path := "./image/default"
	os.Remove(fmt.Sprintf("%s/%s", path, image.Name))

	return s.repo.Delete(userId, itemId)
}
