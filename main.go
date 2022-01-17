package main

import (
	"github.com/HiChen85/RedditRecipesWithGin/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	engin := gin.Default()
	
	engin.GET("/", handlers.HelloWorldGin)
	
	engin.Run("localhost:8000")
}
