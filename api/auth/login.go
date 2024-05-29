package auth

import (
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/result"
	"golang-pet-api/common/utils"
	"golang-pet-api/common/utils/jwt"
	"golang-pet-api/models/model"
	"golang-pet-api/service/user_srv"
)

var userService user_srv.UserService
var redisStore utils.RedisStore

// 登录
func Login(c *gin.Context) {
	var from model.RequestFrom
	c.ShouldBindJSON(&from)
	err := from.Validate()
	if err != nil {
		result.Failed(c, 501, err.Error())
		return
	}
	// 判断验证码是否正确
	verified := redisStore.Verify(from.IdKey, from.Image, false)
	if !verified {
		result.Failed(c, 501, "验证码错误！")
		return
	}
	// 删除验证码缓存
	redisStore.Get(from.IdKey, true)
	// 查询用户
	user := userService.FindUserByLoginFrom(from)
	if user.ID == 0 {
		result.Failed(c, 501, "用户名或密码错误！")
		return
	}
	// 生成token
	token, err := jwt.CreateTokenByUser(user)
	if err != nil {
		result.Failed(c, 501, "token签发失败！")
		return
	}
	jwtUser := model.ConvertToJwtUser(user, "Bearer "+token)
	result.Success(c, jwtUser)
	return
}
