package service

import (
	"github.com/dhevve/uploadImage/internal/models"
	"github.com/dhevve/uploadImage/internal/repository"
)

type UploadService struct {
	repo repository.UploadImage
}

func NewUploadService(repo repository.UploadImage) *UploadService {
	return &UploadService{repo: repo}
}

func (s *UploadService) Upload(userId int, fullFileName string) (int, error) {
	return s.repo.Upload(userId, fullFileName)
}

func (s *UploadService) GetAll(userId int) ([]models.Image, error) {
	return s.repo.GetAll(userId)
}

func (s *UploadService) GetById(userId, imageId int) (models.Image, error) {
	return s.repo.GetById(userId, imageId)
}

func (s *UploadService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}
