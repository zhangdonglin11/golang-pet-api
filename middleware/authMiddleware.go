// 鉴权中间件
package middleware

import (
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/constant"
	"golang-pet-api/common/result"
	"golang-pet-api/common/utils/jwt"
	"strings"
)

func AuthMinddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("token")
		if authHeader == "" {
			result.Failed(c, int(result.ApiCode.NOAUTH), result.ApiCode.GetMessage(result.ApiCode.NOAUTH))
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			result.Failed(c, int(result.ApiCode.AUTHFORMATEROR), result.ApiCode.GetMessage(result.ApiCode.AUTHFORMATEROR))
			c.Abort()
			return
		}
		// 验证token
		jwtUser, err := jwt.ValidateJWT(parts[1])
		if err != nil {
			result.Failed(c, 501, "token验证错误")
			c.Abort()
			return
		}
		c.Set(constant.ContextKeyUserObj, jwtUser)
		c.Set("uid", jwtUser.UserId)
		c.Next()
	}
}
