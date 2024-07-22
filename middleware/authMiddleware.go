// 鉴权中间件
package middleware

import (
	"github.com/gin-gonic/gin"
	"golang-pet-api/utils/jwt"
	"net/http"
	"strings"
)

func AuthOptional() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("token")
		if authHeader == "" {
			c.Next()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "token令牌格式错误",
			})
			c.Abort()
			return
		}
		// 验证token
		claims, err := jwt.ParseJwtToken(jwt.TokenKey, parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "token令牌认证失败",
			})
			c.Abort()
			return
		}
		//将token值解析保存在上下文中
		for k, v := range claims {
			c.Set(k, v)
		}

		c.Next()
	}
}

func AuthMinddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("token")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "token令牌不能为空",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "token令牌格式错误",
			})
			c.Abort()
			return
		}
		// 验证token
		claims, err := jwt.ParseJwtToken(jwt.TokenKey, parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "token令牌认证失败",
			})
			c.Abort()
			return
		}
		//将token值解析保存在上下文中
		for k, v := range claims {
			c.Set(k, v)
		}
		c.Next()
	}
}
