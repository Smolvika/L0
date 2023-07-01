package handler

import (
	"github.com/gin-gonic/gin"
	"orders/pkg/server/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	order := router.Group("/order", h.recovery)

	{
		order.GET("/:id", h.GetOrderById)
	}

	router.LoadHTMLGlob("frontend/view/*")
	router.Static("/static/", "./frontend/static")

	return router
}
