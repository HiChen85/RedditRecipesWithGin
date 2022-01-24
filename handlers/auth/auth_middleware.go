package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
)

func (a *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		claims := jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			// 如果 token 过期,可以直接让其跳转回登录页面
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if token == nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// token 合法且未过期时, 检查当前 token 是否在 redis 黑名单中,如果存在表示之前登出过,不允许再通过验证
		result, _ := a.redisClient.Exists(a.ctx, tokenString).Result()
		// 查到就返回 1
		if result == 1 {
			// 使用 gin 框架的 Abort 函数,会阻止后续中间件函数的执行,同时,相对应的页面 handler 函数也不会执行.
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// 在当前中间件生效之后,去检查后续是否有中间件未处理
		c.Next()
	}
}
