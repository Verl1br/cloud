package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dhevve/uploadImage/internal/models"
	"github.com/gin-gonic/gin"
)

const IMAGE_DIR = "./image"

func (h *Handler) uploadImage(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	imgId, err := h.services.UploadImage.Upload(c, userId)
	if err != nil {
		return
	}
	fmt.Println(imgId)
	c.JSON(http.StatusOK, "ok")
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

	h.services.UploadImage.Delete(userId, id)

	c.JSON(http.StatusOK, "ok")
}
