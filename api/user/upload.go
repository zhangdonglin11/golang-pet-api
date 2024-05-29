package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/constant"
	"golang-pet-api/common/global"
	"golang-pet-api/common/result"
	"golang-pet-api/common/utils/imageUtils"
	"golang-pet-api/models/model"
	"path/filepath"
)

func Upload(c *gin.Context) {
	// 单文件
	file, _ := c.FormFile("file")
	fmt.Println(file.Filename)

	//1.新的文件名
	newFileName := imageUtils.GenerateUniqueFileName() + filepath.Ext(file.Filename)
	//2.保存路径
	newFilePath := "." + global.Config.ImageSettings.UploadDir + newFileName
	err := c.SaveUploadedFile(file, newFilePath)
	if err != nil {
		result.Failed(c, 501, "头像图片保存失败！")
		return
	}
	// 获取token的id
	tokenjwt, _ := c.Get(constant.ContextKeyUserObj)
	jwtUser, _ := tokenjwt.(*model.JwtUser)

	var userChange model.UserChange
	userChange.UserID = jwtUser.UserId
	// 查询用户的修改表
	global.Db.Take(&userChange)
	// 删除旧的图片
	imageUtils.DeleteImage(userChange.Avatar)

	//保存新的头像图片名
	userChange.Avatar = global.Config.ImageSettings.UploadDir + newFileName
	userChange.Approved = true
	global.Db.Save(&userChange)
	global.Db.Take(&userChange)
	userChange.Avatar = global.Config.ImageSettings.ImageHost + userChange.Avatar
	result.Success(c, userChange)
	return
}
