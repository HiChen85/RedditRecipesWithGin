package handlers

import (
	"context"
	"encoding/json"
	"github.com/HiChen85/RedditRecipesWithGin/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type RecipeHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewRecipeHandler(ctx context.Context, mongoCollection *mongo.Collection, redisClient *redis.Client) *RecipeHandler {
	return &RecipeHandler{
		collection:  mongoCollection,
		ctx:         ctx,
		redisClient: redisClient,
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while insert one",
		})
		return
	}
	log.Println("Successfully insert one to MongoDB...")
	
	_, err = r.redisClient.Get(r.ctx, "recipes").Result()
	if err != redis.Nil {
		log.Println("Delete cache from Redis...")
		r.redisClient.Del(r.ctx, "recipes")
	} else if err != nil { // 保证程序的健壮性
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func (r *RecipeHandler) ListRecipesHandler(c *gin.Context) {
	// 在向 mongodb 检索数据之前, 先查询缓存中是否有数据
	// 当 redis 中不存在检索的 key, 那么返回值是一个 redis.Nil 类型
	val, err := r.redisClient.Get(r.ctx, "recipes").Result()
	
	// 当 redis 中没有对应缓存, 就向 mongodb 中检索
	if err == redis.Nil {
		log.Println("Request to mongoDB....")
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
		
		// 成功返回后, 再向 redis 中添加对应缓存
		cacheData, _ := json.Marshal(recipes)
		// 最后一个参数代表缓存过期时间, 0 代表缓存永不过期
		// redis 中的数据类型有 String, 所以这里可以直接将所需要的数据转换成 string
		r.redisClient.Set(r.ctx, "recipes", string(cacheData), 0)
		c.JSON(http.StatusOK, recipes)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Println("Request to Redis....")
		recipes := make([]*models.Recipe, 0)
		err := json.Unmarshal([]byte(val), &recipes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, recipes)
	}
	
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
	log.Println("successfully Update...")
	_, err = r.redisClient.Get(r.ctx, "recipes").Result()
	if err != redis.Nil {
		log.Println("Delete cache from Redis")
		r.redisClient.Del(r.ctx, "recipes")
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
		return
	}
	one, err := r.collection.DeleteOne(r.ctx, bson.M{"_id": objID})
	// 成功删除已存在的数据是, 返回的 DeletedCount 结果就是 1
	if one.DeletedCount != 1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "no such data matched",
		})
		return
	}
	log.Println("Successfully delete")
	_, err = r.redisClient.Get(r.ctx, "recipes").Result()
	if err != redis.Nil {
		log.Println("Delete cache from Redis")
		r.redisClient.Del(r.ctx, "recipes")
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully delete document",
	})
}

func (r *RecipeHandler) SearchRecipeHandler(c *gin.Context) {
	val, err := r.redisClient.Get(r.ctx, "recipesWithTags").Result()
	if err == redis.Nil {
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
		// 因为 redis 中存入的是 string类型,所以必须先将 go 的结构体对象转为 byte 数组,再转为 string
		cacheData, _ := json.Marshal(tempRecipes)
		r.redisClient.Set(r.ctx, "recipesWithTags", string(cacheData), 0)
		c.JSON(http.StatusOK, tempRecipes)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Println("Request to Redis...")
		recipes := make([]*models.Recipe, 0)
		err := json.Unmarshal([]byte(val), &recipes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, recipes)
	}
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
