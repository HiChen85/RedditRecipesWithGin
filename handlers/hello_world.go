package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelloWorldGin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}
