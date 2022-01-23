package handlers

import (
	"context"
	"github.com/HiChen85/RedditRecipesWithGin/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type RecipeHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewRecipeHandler(collection *mongo.Collection, ctx context.Context) *RecipeHandler {
	return &RecipeHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (r *RecipeHandler) PostNewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishAt = time.Now()
	_, err := r.collection.InsertOne(r.ctx, &recipe)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, recipe)
}

func (r *RecipeHandler) ListRecipesHandler(c *gin.Context) {
	// Find 函数返回的 cursor 位于第一个数据之前,第一次调用 Next 后会指向第一个数据
	cur, err := r.collection.Find(r.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer cur.Close(r.ctx)
	recipes := make([]*models.Recipe, 0)
	for cur.Next(r.ctx) {
		var currentData models.Recipe
		// 将数据临时解析到一个结构体上,然后将数据写入到提前定义好的食谱数组中
		err := cur.Decode(&currentData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			log.Fatal(err)
		}
		recipes = append(recipes, &currentData)
	}
	c.JSON(http.StatusOK, recipes)
}

func (r *RecipeHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	objectID, _ := primitive.ObjectIDFromHex(id)
	// bson.M 用于提供条件查询过滤, bson.D 用于更新, 更新时使用 $set 参数为 key,表示更新文档
	// 再传入一个 bson.D 作为数据更新部分
	_, err := r.collection.UpdateOne(r.ctx, bson.M{"_id": objectID}, bson.D{
		{"$set", bson.D{
			{"name", recipe.Name},
			{"instructions", recipe.Instructions},
			{"ingredients", recipe.Ingredients},
			{"tags", recipe.Tags},
		}},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, recipe)
}

func (r *RecipeHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	// 利用此参数来获取数据库中存储的文档 id
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		log.Fatal(err)
	}
	one, err := r.collection.DeleteOne(r.ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		log.Fatal(err)
	}
	log.Println("删除文档数:", one.DeletedCount)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully delete document",
	})
}

func (r *RecipeHandler) SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	tempRecipes := make([]*models.Recipe, 0)
	// 使用 bson.M 作为过滤器, 当进行 in 查询时, 先设置好要查询的字段, 然后重新设置一个 bson.M 设置 in 参数
	// 在设定好对应字段中匹配的值
	cursor, err := r.collection.Find(r.ctx, bson.M{"tags": bson.M{"$in": []string{tag}}})
	defer cursor.Close(r.ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		log.Fatal(err)
	}
	if err = cursor.All(r.ctx, &tempRecipes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, tempRecipes)
}

func (r *RecipeHandler) GetOneRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(id)
	one := r.collection.FindOne(r.ctx, bson.M{"_id": objID})
	var tempData models.Recipe
	err := one.Decode(&tempData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, tempData)
}
