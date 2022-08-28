package handler

import (
	"github.com/dhevve/uploadImage/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Static("/stat-img", "./image")

	auth := router.Group("/auth")
	{
		auth.POST("/sing-up", h.signUp)
		auth.POST("/sing-in", h.signIn)
	}

	upload := router.Group("/api", h.userIdentity)
	{
		upload.POST("/upload_image", h.uploadImage)
		upload.GET("/", h.getAll)
		upload.GET("/:id", h.getById)
		upload.GET("/download/:id", h.download)
		upload.DELETE("/:id", h.deleteImage)
	}

	return router
}
