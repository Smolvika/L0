package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) GetOrderById(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "badRequest.html", nil)
		return
	}

	order, err := h.services.GetOrderById(orderId)
	if err != nil {
		c.HTML(http.StatusNotFound, "notFound.html", nil)
		return
	}

	c.HTML(http.StatusOK, "index.html", order)
}
