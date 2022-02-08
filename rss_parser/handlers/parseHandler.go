package handlers

import (
	"context"
	"github.com/HiChen85/RedditRecipesWithGin/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
)

var mongoClient *mongo.Client
var ctx context.Context

func init() {
	ctx = context.Background()
	Client, err := mongo.Connect(ctx, options.Client().ApplyURI(utils.MONGO_URI))
	mongoClient = Client
	if err != nil {
		log.Fatal(err)
	}
	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Rss Parser Successfully Connects to MongoDB...")
}

type requestData struct {
	URL string `json:"url"`
}

// ParseHandler 获取食谱数据后将数据插入数据库
func ParseHandler(c *gin.Context) {
	// 根据数据库连接拿到或者创建一个新的mongodb 集合
	recipeCollection := mongoClient.Database(utils.MONGO_DATABASE).Collection(utils.MONGO_RSS_PARSER_COLLECTION)
	
	jsonData := new(requestData)
	if err := c.ShouldBindJSON(jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 成功解析到上传的 JSON 数据后调用获取数据的方法, 从reddit 中获取数据
	entries, err := utils.GetDataFromReddit(jsonData.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法获取数据" + err.Error(),
		})
		return
	}
	
	// 根据分析 reddit 的网站,前两个帖子的内容并非所需要的数据,通过切片的方式过滤
	for i := range entries[2:] {
		// InsertOne 会返回一个插入结果,成功插入的话,可以从 InsertOne 对象中获取到插入完成后的 ID
		_, err := recipeCollection.InsertOne(ctx, bson.M{
			"title": entries[i].Title,
			// 帖子配图地址
			"thumbnail": entries[i].Thumbnail.URL,
			// 帖子地址
			"link": entries[i].Link.Href,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "成功插入", "data": entries})
}
