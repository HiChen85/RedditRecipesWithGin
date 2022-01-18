package main

import (
	"github.com/HiChen85/RedditRecipesWithGin/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	
	// this is an init router for gin
	engine.GET("/", handlers.HelloWorldGin)
	
	// routers for Recipe
	engine.POST("/recipes", handlers.NewRecipeHandler)
	engine.GET("/recipes", handlers.ListRecipesHandler)
	engine.PUT("/recipes/:id", handlers.UpdateRecipeHandler)
	
	engine.Run("localhost:8000")
}
