package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-pet-api/service/avatar_srv"
	"net/http"
)

// Captcha 获取验证码接口
func Captcha(c *gin.Context) {
	id, base64Image, err := avatar_srv.CaptMake()
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    1,
		"message": "请求成功",
		"data": map[string]interface{}{
			"idKey":      id,
			"imageUtils": base64Image,
		},
	})
}
