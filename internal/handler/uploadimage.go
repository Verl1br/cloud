package handler

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dhevve/uploadImage/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid"
)

const IMAGE_DIR = "./image"

func (h *Handler) uploadImage(c *gin.Context) {
	productId, _ := c.GetPostForm("product_id")
	{
		if len(productId) == 0 {
			productId = "default"
		}
		err := os.MkdirAll(fmt.Sprintf("%s/product/%s", IMAGE_DIR, productId), os.ModePerm) // создаем директорию с id продукта, если еще не создана
		if err != nil {
			return
		}
	}
	path := fmt.Sprintf("%s/product/%s", IMAGE_DIR, productId)

	// извлекаем файл из парамeтров post запроса
	form, _ := c.MultipartForm()
	var fileName string
	imgExt := "jpeg"
	// берем первое имя файла из присланного списка
	for key := range form.File {
		fileName = key
		// извлекаем расширение файла
		arr := strings.Split(fileName, ".")
		if len(arr) > 1 {
			imgExt = arr[len(arr)-1]
		}
		continue
	}

	// извлекаем содержание присланного файла по названию файла
	file, _, err := c.Request.FormFile(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	// читаем содержание присланного файл в []byte
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	fullFileName := fmt.Sprintf("%s.%s", randomFilename(), imgExt)

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	imgId, err := h.services.UploadImage.Upload(userId, fullFileName)
	if err != nil {
		return
	}
	// открываем файл для сохранения картинки
	fileOnDisk, err := os.Create(fmt.Sprintf("%s/%s", path, fullFileName))
	if err != nil {
		return
	}
	defer fileOnDisk.Close()

	_, err = fileOnDisk.Write(fileBytes)
	if err != nil {
		return
	}

	fmt.Println(imgId)
	c.JSON(http.StatusOK, "ok")
}

func randomFilename() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return strings.ToLower(fmt.Sprintf("%v", ulid.MustNew(ulid.Timestamp(t), entropy)))
}

type getAllResponse struct {
	Data []models.Image `json:"data"`
}

func (h *Handler) getAll(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	images, err := h.services.UploadImage.GetAll(userId)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, getAllResponse{
		Data: images,
	})
}

func (h *Handler) getById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}

	image, err := h.services.UploadImage.GetById(userId, id)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, image)
}

func (h *Handler) deleteImage(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}

	h.deleteFile(userId, id)

	h.services.UploadImage.Delete(userId, id)

	c.JSON(http.StatusOK, "ok")
}

func (h *Handler) deleteFile(userId, id int) {
	image, err := h.services.UploadImage.GetById(userId, id)
	if err != nil {
		return
	}
	path := "./image/product/default"
	os.Remove(fmt.Sprintf("%s/%s", path, image.Name))
}
