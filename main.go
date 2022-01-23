package main

import (
	"context"
	"github.com/HiChen85/RedditRecipesWithGin/handlers"
	"github.com/HiChen85/RedditRecipesWithGin/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var recipeHandler *handlers.RecipeHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(utils.MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB......")
	mongoCollection := client.Database(utils.MONGO_DATABASE).Collection(utils.MONGO_COLLECTION)
	redisClient := redis.NewClient(&utils.RedisOptions)
	status := redisClient.Ping(ctx)
	log.Println("Redis status:", status)
	recipeHandler = handlers.NewRecipeHandler(ctx, mongoCollection, redisClient)
}

func main() {
	engine := gin.Default()
	
	// this is an init router for gin
	engine.GET("/", handlers.HelloWorldGin)
	
	// routers group for Recipe
	recipes := engine.Group("/recipes")
	{
		recipes.POST("/", recipeHandler.PostNewRecipeHandler)
		recipes.GET("/", recipeHandler.ListRecipesHandler)
		recipes.PUT("/:id", recipeHandler.UpdateRecipeHandler)
		recipes.DELETE("/:id", recipeHandler.DeleteRecipeHandler)
		recipes.GET("/search", recipeHandler.SearchRecipeHandler)
		recipes.GET("/:id", recipeHandler.GetOneRecipeHandler)
	}
	
	engine.Run("localhost:8000")
}
