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
	recipes := engine.Group("/recipes")
	{
		recipes.POST("/", handlers.NewRecipeHandler)
		recipes.GET("/", handlers.ListRecipesHandler)
		recipes.PUT("/:id", handlers.UpdateRecipeHandler)
		recipes.DELETE("/:id", handlers.DeleteRecipeHandler)
	}
	
	engine.Run("localhost:8000")
}
