package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/constant"
	"golang-pet-api/common/global"
	"golang-pet-api/common/result"
	"golang-pet-api/models/model"
)

func GetUser(c *gin.Context) {
	tokenjwt, _ := c.Get(constant.ContextKeyUserObj)
	jwtUser, ok := tokenjwt.(*model.JwtUser)
	if !ok {
		// 类型断言失败，处理错误
		fmt.Println("无法将 tokenjwt 转换为 *jwtUser 类型")
		return
	}
	var user model.User
	global.Db.Take(&user, jwtUser.UserId)
	if user.ID == 0 {
		result.Failed(c, 501, "没有查询到该用户")
	}
	newjwtUser := model.ConvertToJwtUser(user, "")
	result.Success(c, newjwtUser)
	return
}
