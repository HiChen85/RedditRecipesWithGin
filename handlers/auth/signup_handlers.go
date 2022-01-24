package auth

import (
	"crypto/sha256"
	"github.com/HiChen85/RedditRecipesWithGin/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
	"time"
)

func (a *AuthHandler) SignUpHandler(c *gin.Context) {
	// 注册功能:
	// 前端校验两次输入的密码一致后, 将用户名与密码传到后台,后台验证该用户是否存在,若存在则返回,不存在则创建新用户
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 成功绑定了用户对象后, 利用这个对象向数据库检索
	h := sha256.New()
	result := a.collection.FindOne(a.ctx, bson.M{
		"username": newUser.Username,
		"password": string(h.Sum([]byte(newUser.Password))),
	})
	// 当 Err 为空时,表示数据库中存在这个用户
	if result.Err() == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "this account has already existed",
		})
		return
	}
	_, err := a.collection.InsertOne(a.ctx, bson.M{
		"username": newUser.Username,
		"password": string(h.Sum([]byte(newUser.Password))),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Created New Account",
	})
}

func (a *AuthHandler) SignOutHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	
	// 此处 err 不为空有以下几种情况:
	// 1. 当不携带 Authorization 时,
	// 2. 当 token 过期时
	// 假设在登陆之后,所有的请求都会携带 token, 所以暂时忽略第一种情况, 这里仅仅处理 token 过期
ExpireError:
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "please login again",
		})
		return
	}
	if token == nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}
	// 如果 token 没过期时登出, 将当前 token 剩下的过期时间算出来, 并将当前的 token 加入到 redis 的黑名单中
	duration := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
	if duration < 30*time.Second {
		goto ExpireError
	}
	// duration 大于 30 秒, 表示未过期, 则将当前 token 加入到 redis 中表示登出
	ok, _ := a.redisClient.SetNX(a.ctx, tokenString, "", duration).Result()
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "add blacklist error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "SignOut",
	})
}
