package main

import (
	"context"
	"embed"
	"github.com/HiChen85/RedditRecipesWithGin/handlers"
	"github.com/HiChen85/RedditRecipesWithGin/handlers/auth"
	"github.com/HiChen85/RedditRecipesWithGin/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed assets/* templates/*
var FS embed.FS

var authHandler *auth.AuthHandler
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
	authHandler = auth.NewAuthHandler(ctx, authCollection, redisClient)
}

func main() {
	
	// 在 ParseFS 的第二个参数 patterns 中, 可以设置多个匹配模式, 用于匹配不同类型的文件,但这些文件后缀必须存在于目录下,如果不存在则报错
	tmpl := template.Must(template.New("").ParseFS(FS, "templates/*.html"))
	fsAssets, err := fs.Sub(FS, "assets")
	if err != nil {
		panic(err)
	}
	
	// 定义 gin 路由
	engine := gin.Default()
	
	// 定义 cors 中间件
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	// 通过 go embed 命令,导入目录下的静态资源和模板文件用于打包.
	engine.SetHTMLTemplate(tmpl)
	engine.StaticFS("/assets", http.FS(fsAssets))
	
	// this is an init router for gin
	engine.GET("/", handlers.HelloWorldGin)
	engine.POST("/signin", authHandler.SignInHandler)
	engine.POST("/signup", authHandler.SignUpHandler)
	engine.POST("/refresh", authHandler.RefreshTokenHandler)
	engine.POST("/signout", authHandler.SignOutHandler)
	engine.GET("/recipes", recipeHandler.ListRecipesHandler)
	
	// routers group for Recipe
	recipes := engine.Group("/recipes")
	// 使用认证中间件来进行用户验证.
	recipes.Use(authHandler.AuthMiddleware())
	{
		//recipes.GET("/", recipeHandler.ListRecipesHandler)
		recipes.POST("/", recipeHandler.PostNewRecipeHandler)
		recipes.PUT("/:id", recipeHandler.UpdateRecipeHandler)
		recipes.GET("/:id", recipeHandler.GetOneRecipeHandler)
		recipes.DELETE("/:id", recipeHandler.DeleteRecipeHandler)
		recipes.GET("/search", recipeHandler.SearchRecipeHandler)
	}
	
	engine.Run("localhost:8000")
}
