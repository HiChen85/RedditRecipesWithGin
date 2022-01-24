package main

import (
	"context"
	"github.com/HiChen85/RedditRecipesWithGin/handlers"
	"github.com/HiChen85/RedditRecipesWithGin/handlers/middlewares"
	"github.com/HiChen85/RedditRecipesWithGin/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var authHandler *handlers.AuthHandler
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
	// create mongodb recipes collection
	recipesCollection := client.Database(utils.MONGO_DATABASE).Collection(utils.MONGO_RECIPES_COLLECTION)
	authCollection := client.Database(utils.MONGO_DATABASE).Collection(utils.MONGO_USER_COLLECTION)
	// New Redis Client
	redisClient := redis.NewClient(&utils.RedisOptions)
	status := redisClient.Ping(ctx)
	log.Println("Redis status:", status)
	// recipeHandler
	recipeHandler = handlers.NewRecipeHandler(ctx, recipesCollection, redisClient)
	// authHandler
	authHandler = handlers.NewAuthHandler(ctx, authCollection, redisClient)
}

func main() {
	engine := gin.Default()
	
	// this is an init router for gin
	engine.GET("/", handlers.HelloWorldGin)
	engine.POST("/signin", authHandler.SignInHandler)
	engine.POST("/refresh", authHandler.RefreshTokenHandler)
	
	// routers group for Recipe
	recipes := engine.Group("/recipes")
	// 使用认证中间件来进行用户验证.
	recipes.Use(middlewares.AuthMiddleware())
	{
		recipes.GET("/", recipeHandler.ListRecipesHandler)
		recipes.POST("/", recipeHandler.PostNewRecipeHandler)
		recipes.PUT("/:id", recipeHandler.UpdateRecipeHandler)
		recipes.GET("/:id", recipeHandler.GetOneRecipeHandler)
		recipes.DELETE("/:id", recipeHandler.DeleteRecipeHandler)
		recipes.GET("/search", recipeHandler.SearchRecipeHandler)
	}
	
	engine.Run("localhost:8000")
}
