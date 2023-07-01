package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) recovery(c *gin.Context) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("There was an internal server error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, "There was an internal server error")
		}
	}()
}
