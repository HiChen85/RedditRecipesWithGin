package handlers

import (
	"context"
	"github.com/HiChen85/RedditRecipesWithGin/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

type RecipeData struct {
	Title     string `bson:"title"`
	Link      string `bson:"link"`
	Thumbnail string `bson:"thumbnail"`
}

func DashboardHandler(c *gin.Context) {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(utils.MONGO_URI))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	collection := mongoClient.Database(utils.MONGO_DATABASE).Collection(utils.MONGO_CONSUMER_COLLECTION)
	// 获取所有的数据,展示到页面上
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipes := make([]*RecipeData, 0)
	
	for cur.Next(context.TODO()) {
		var tempData RecipeData
		cur.Decode(&tempData)
		recipes = append(recipes, &tempData)
	}
	// 将 recipes 渲染回页面上, 利用模板引擎
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"recipes": recipes,
	})
	
}
