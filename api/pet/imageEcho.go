package pet

import (
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/global"
	"golang-pet-api/common/result"
	"golang-pet-api/common/utils/imageUtils"
	"golang-pet-api/common/utils/rabbitMQ"
	"net/http"
	"path/filepath"
)

// ImageEcho 上传图片回显
func ImageEcho(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 生产文件名
	newFileName := imageUtils.GenerateUniqueFileName() + filepath.Ext(file.Filename)
	// 保存地址+文件名
	uploadDir := global.Config.ImageSettings.UploadDir + newFileName
	newFilePath := filepath.Join(".", uploadDir)
	if err = c.SaveUploadedFile(file, newFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 将不确定是否永久保存图片信息添加到redis中,定期删除临时图片 回显上传的图片
	err = imageUtils.TempRedisStore{}.Set(newFileName, uploadDir)
	if err != nil {
		return
	}

	rabbitMQ.SandTempImg(newFileName)

	result.Success(c, map[string]string{
		"imageUrl": global.Config.ImageSettings.ImageHost + uploadDir,
	})
}
