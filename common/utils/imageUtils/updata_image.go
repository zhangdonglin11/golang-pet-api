package imageUtils

import (
	"errors"
	"fmt"
	"golang-pet-api/common/global"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// ResetAvatar 重置图片
func ResetAvatar() (string, error) {
	// 保存图片的方法
	rand.Seed(time.Now().UnixNano())
	nub := rand.Intn(4) + 1
	// 源图片路径
	sourceImagePath := fmt.Sprintf("./static/avatar/ikun%d.png", nub)
	// 打开源图片文件
	srcFile, err := os.Open(sourceImagePath)
	if err != nil {
		return "", errors.New("avatar打开文件失败")
	}
	defer srcFile.Close()

	// 目标文件夹路径
	targetDir := "." + global.Config.ImageSettings.UploadDir
	// 创建目标文件夹（如果不存在）
	if err = os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", errors.New("创建目标文件夹失败")
	}
	// 提取源图片文件名
	_, fileName := filepath.Split(sourceImagePath)
	// 获取文件后缀名
	fileExt := filepath.Ext(fileName)
	// 新文件名
	newFileName := GenerateUniqueFileName() + fileExt
	// 创建目标图片文件
	targetImagePath := filepath.Join(targetDir, newFileName)
	dstFile, err := os.Create(targetImagePath)
	if err != nil {
		return "", errors.New("创建avatar文件失败")
	}
	defer dstFile.Close()
	// 将源图片内容复制到目标图片文件
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return "", errors.New("复制文件失败")
	}
	return global.Config.ImageSettings.UploadDir + newFileName, nil
}

// DeleteImage 删除图片的方法
func DeleteImage(filePath string) bool {
	// 图片完整路径 + 图片文件名
	filePath = "." + filePath
	// 检查文件是否存在
	if _, err := os.Stat(filePath); err == nil {
		// 文件存在，可以删除
		if err := os.Remove(filePath); err != nil {
			fmt.Println("删除文件失败:", err)
			return false
		}
		fmt.Println("文件删除成功")
		return true
	} else if os.IsNotExist(err) {
		// 文件不存在
		fmt.Println("文件不存在")
		return false
	} else {
		// 其他错误
		fmt.Println("检查文件状态时出错:", err)
		return false
	}
}

// GenerateUniqueFileName  生成唯一的文件名
func GenerateUniqueFileName() string {
	// 使用当前时间戳作为文件名的一部分
	timestamp := time.Now().UnixNano()
	// 生成一个随机字符串作为文件名的一部分
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	randomStr := make([]rune, 8)
	for i := range randomStr {
		randomStr[i] = letters[rand.Intn(len(letters))]
	}
	// 将时间戳和随机字符串拼接成文件名
	fileName := fmt.Sprintf("%d_%s", timestamp, string(randomStr))
	return fileName
}
