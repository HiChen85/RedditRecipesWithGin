package main

// 用来处理数据获取的爬虫程序

import (
	"github.com/HiChen85/RedditRecipesWithGin/rss_parser/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	// 利用 POST 方法上传后缀为.rss 的 url, 交给 ParseHandler 处理
	engine.POST("/parse", handlers.ParseHandler)
	engine.GET("/recipes", handlers.GetAllRecipes)
	engine.Run(":8001")
}
