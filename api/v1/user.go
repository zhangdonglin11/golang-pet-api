package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-pet-api/global"
	"golang-pet-api/models/forms"
	"golang-pet-api/models/model"
	"golang-pet-api/utils"
	"golang-pet-api/utils/imageUtils"
	"path/filepath"
)

type User struct {
}

// GetUser godoc
// @Summary 获取用户详细信息
// @Description 通过请求头附带 Bearer Token 获取用户详细信息。
// @Tags 用户模块
// @Produce json
// @Security ApiKeyAuth
// @Param token header string true "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Success 200 {string} json{Code,Msg,Data}  "成功"
// @Router /api/v1/user [get]
func (u User) GetUser(c *gin.Context) {
	fmt.Println(c.Get("userId"))
	f64 := c.GetFloat64("userId")
	userId := uint(f64)
	var user model.User
	if result := global.Db.Preload("Profile").First(&user, userId); result.RowsAffected == 0 {
		utils.RespFail(c.Writer, "用户不存在!")
		return
	}
	user.Profile.Icon = global.Config.ImageSettings.ImageHost + user.Profile.Icon
	utils.RespOk(c.Writer, user.Profile, "请求成功")
}

// UpdateUser godoc
// @Summary 修改用户信息
// @Tags 用户模块
// @Produce json
// @Security ApiKeyAuth
// @Param token header string true "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param nickname formData string false "昵称"
// @Param icon     formData string false "头像url"
// @Param phone    formData string false "电话"
// @Param email    formData string false "邮箱"
// @Param gender   formData string false "性别"
// @Param status   formData string false "状态"
// @Success 200 {string} json{Code,Msg,Data}  "成功"
// @Router /api/v1/user [put]
func (User) UpdateUser(c *gin.Context) {
	f64 := c.GetFloat64("userId")
	userId := uint(f64)

	// 绑定并验证上传的数据
	var form forms.ProfileForm
	if err := c.ShouldBind(&form); err != nil {
		HandleValidatorError(c, err)
		return
	}
	var profile model.Profile
	if result := global.Db.Where("user_id=?", userId).First(&profile); result.RowsAffected == 0 {
		utils.RespFail(c.Writer, "用户信息查询错误")
		return
	}
	profile.Nickname = form.Nickname
	profile.Icon = form.Icon
	profile.Phone = form.Phone
	profile.Email = form.Email
	profile.Gender = form.Gender
	profile.Status = form.Status

	if result := global.Db.Model(&profile).Updates(profile); result.RowsAffected == 0 {
		utils.RespFail(c.Writer, "修改个人信息错误")
		return
	}
	utils.RespOk(c.Writer, profile, "修改个人信息成功")
	return
}

// UploadPhoto godoc
// @Summary 上传用户头像
// @Tags 用户模块
// @Produce json
// @Security ApiKeyAuth
// @Param token header string true "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param file formData file true "上传头像图片"
// @Success 200 {string} json{Code,Msg,Data}  "成功"
// @Router /api/v1/user/upload [put]
func (User) UploadPhoto(c *gin.Context) {
	f64 := c.GetFloat64("userId")
	userId := uint(f64)

	f, _ := c.Get("userId")
	switch f.(type) {
	case int:

	}

	// 单文件
	file, _ := c.FormFile("file")
	fmt.Println(file.Filename)

	//1.新的文件名
	newFileName := imageUtils.GenerateUniqueFileName() + filepath.Ext(file.Filename)
	//2.保存路径
	newFilePath := "." + global.Config.ImageSettings.UploadDir + newFileName
	err := c.SaveUploadedFile(file, newFilePath)
	if err != nil {
		utils.RespFail(c.Writer, err.Error())
		return
	}

	var profile model.Profile
	if result := global.Db.First(&profile, "user_id=?", userId); result.RowsAffected == 0 {
		utils.RespFail(c.Writer, "用户不存在")
		return
	}

	// 删除旧的图片
	imageUtils.DeleteImage(profile.Icon)

	//保存新的头像图片名
	profile.Icon = global.Config.ImageSettings.UploadDir + newFileName
	if result := global.Db.Save(&profile); result.RowsAffected == 0 {
		utils.RespFail(c.Writer, "修改用户头像失败")
		return
	}
	utils.RespOk(c.Writer, profile, "修改用户头像成功")
	return
}
