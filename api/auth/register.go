package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/result"
	"golang-pet-api/common/utils/imageUtils"
	"golang-pet-api/models/model"
	"math/rand"
	"time"
)

func Register(c *gin.Context) {
	var from model.RequestFrom
	c.ShouldBindJSON(&from)
	// 判断验证码是否正确
	if verified := redisStore.Verify(from.IdKey, from.Image, false); !verified {
		result.Failed(c, 501, "验证码错误！")
		return
	}
	// 删除验证码缓存
	redisStore.Get(from.IdKey, true)

	// 查询用户是否已经注册
	user := userService.FindUserByUserName(from)
	if user.ID != 0 {
		result.Failed(c, 501, "用户名已被注册")
	}

	//初始化用户表
	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("铲屎官%d", rand.Intn(9000)+1000)
	// 生成初始头像图片
	avatar, _ := imageUtils.ResetAvatar()
	newUser := model.User{
		Username: from.Username,
		Password: from.Password,
		NickName: name,
		Avatar:   avatar,
		Role:     0,
		Status:   true,
		UserChange: model.UserChange{
			NickName: name,
			Approved: true,
		},
	}
	err := userService.SaveUser(&newUser)
	if err != nil {
		result.Failed(c, 501, "注册失败！")
	}
	result.Success(c, "注册成功！")
	return
}
