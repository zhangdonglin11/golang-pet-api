package user

import (
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/constant"
	"golang-pet-api/common/global"
	"golang-pet-api/common/result"
	"golang-pet-api/models/model"
)

// 修改用户信息
func UpdateUser(c *gin.Context) {
	tokenjwt, _ := c.Get(constant.ContextKeyUserObj)
	jwtUser, _ := tokenjwt.(*model.JwtUser)

	var putUserFrom model.PutUserFrom
	c.ShouldBindJSON(&putUserFrom)
	err := putUserFrom.Validate()
	if err != nil {
		result.Failed(c, 501, err.Error())
		return
	}
	var user model.User
	user.ID = jwtUser.UserId
	global.Db.Preload("UserChange").Take(&user)
	global.Db.Save(&model.UserChange{
		User:     &user,
		ID:       user.UserChange.ID,
		NickName: putUserFrom.NickName,
		Gender:   putUserFrom.Gender,
		Tel:      putUserFrom.Tel,
		Approved: false,
	})
	global.Db.Preload("UserChange").Take(&user)
	result.Success(c, user.UserChange)
}
