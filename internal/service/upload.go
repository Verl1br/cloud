package service

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dhevve/uploadImage/internal/models"
	"github.com/dhevve/uploadImage/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid"
)

const IMAGE_DIR = "./image"

type UploadService struct {
	repo repository.UploadImage
}

func NewUploadService(repo repository.UploadImage) *UploadService {
	return &UploadService{repo: repo}
}

func (s *UploadService) Upload(c *gin.Context, userId int) (int, error) {
	productId, _ := c.GetPostForm("product_id")
	{
		if len(productId) == 0 {
			productId = "default"
		}
		err := os.MkdirAll(fmt.Sprintf("%s/product/%s", IMAGE_DIR, productId), os.ModePerm)
		if err != nil {
			return 0, err
		}
	}
	path := fmt.Sprintf("%s/product/%s", IMAGE_DIR, productId)

	form, _ := c.MultipartForm()
	var fileName string
	imgExt := "jpeg"
	for key := range form.File {
		fileName = key
		arr := strings.Split(fileName, ".")
		if len(arr) > 1 {
			imgExt = arr[len(arr)-1]
		}
		continue
	}

	file, _, err := c.Request.FormFile(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, err
	}

	fullFileName := fmt.Sprintf("%s.%s", randomFilename(), imgExt)
	fileOnDisk, err := os.Create(fmt.Sprintf("%s/%s", path, fullFileName))
	if err != nil {
		return 0, err
	}
	defer fileOnDisk.Close()

	_, err = fileOnDisk.Write(fileBytes)
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
	path := "./image/product/default"
	os.Remove(fmt.Sprintf("%s/%s", path, image.Name))

	return s.repo.Delete(userId, itemId)
}

func randomFilename() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return strings.ToLower(fmt.Sprintf("%v", ulid.MustNew(ulid.Timestamp(t), entropy)))
}
