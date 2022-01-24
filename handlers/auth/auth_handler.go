package auth

import (
	"context"
	"crypto/sha256"
	"github.com/HiChen85/RedditRecipesWithGin/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"time"
)

type AuthHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewAuthHandler(ctx context.Context, mongoCollection *mongo.Collection, redisClient *redis.Client) *AuthHandler {
	return &AuthHandler{
		collection:  mongoCollection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

func (a *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 向数据库验证用户
	h := sha256.New()
	result := a.collection.FindOne(a.ctx, bson.M{
		"username": user.Username,
		"password": string(h.Sum([]byte(user.Password))),
	})
	if result.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}
	// 当验证通过, 为用户生成一个 token, 过期时间为当前时间 + 几分钟
	expirationTime := time.Now().Add(3 * time.Minute)
	// 定义一个标准的 payload 中的部分字段
	claims := jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Issuer:    "Recipes App",
	}
	// 当前生成的只是 token 对象,要生成真正的 token 还需要利用秘钥
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 当前利用设置环境变量的方式手动设置秘钥
	secretString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": secretString,
	})
}

func (a *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	defer func() {
		// 如果是携带的 token 已经过期了,那么会走到这里处理
		if err := recover(); err != nil {
			log.Println(err)
			expirationTime := time.Now().Add(5 * time.Minute)
			claims.ExpiresAt = expirationTime.Unix()
			token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			secretString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"token": secretString,
			})
		}
	}()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		// 如果自带的这个过期时间处理错误产生了,就去 defer 处理
		panic(err)
	}
	if token == nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}
	// token 最终过期的时间减去当前时间大于 30 秒. 认为 token 未过期
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Token is not expired",
		})
		return
	}
	// 如果 token 过期, 那么就重新设置一个新的 token 返回
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": secretString,
	})
}
