package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-pet-api/global"
	"golang-pet-api/utils"
	"golang-pet-api/utils/imageUtils"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

// UploadImageEcho godoc
// @Summary 上传图片+回显
// @Tags 上传图片模块
// @Produce json
// @Security ApiKeyAuth
// @Param token header string true "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param file formData file true "上传图片"
// @Success 200 {string} json{Code,Msg,Data}  "成功"
// @Router /api/v1/upload [post]
func UploadImageEcho(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.RespFail(c.Writer, err.Error())
		return
	}
	// 生产文件名+后缀
	newFileName := imageUtils.GenerateUniqueFileName() + filepath.Ext(file.Filename)
	// 拼接保存地址
	uploadDir := global.Config.ImageSettings.UploadDir + newFileName
	// 拼接个点
	newFilePath := filepath.Join(".", uploadDir)
	if err = c.SaveUploadedFile(file, newFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 将不确定是否永久保存图片信息添加到redis中,定期删除临时图片 回显上传的图片
	// 获取当前时间
	now := time.Now()
	// 生成一个20到40分钟之间的随机偏移量（单位：秒）
	rand.New(rand.NewSource(time.Now().UnixNano()))
	offset := rand.Intn(21) + 10 // 10到30之间的随机数
	// 计算目标到期时间
	targetTime := now.Add(time.Minute * time.Duration(offset))
	data := map[string]interface{}{
		"expirationTime": targetTime.Unix(),
		"uploadDir":      uploadDir,
	}
	//将临时图片的信息保存到redis中
	if err = global.RedisDb.HMSet(c, global.TempFile+newFileName, data).Err(); err != nil {
		utils.RespFail(c.Writer, err.Error())
		return
	}
	utils.RespOk(c.Writer, global.Config.ImageSettings.ImageHost+uploadDir, "临时图片删除时间:"+targetTime.Format("2006-01-02 15:04:05"))
	return
}

// HandleTimeOutFile 处理临时图片 异步
func HandleTimeOutFile(ctx context.Context, FileCh *chan bool) {
	// 初始化 SCAN 游标为 0
	client := global.RedisDb
	var cursor uint64 = 0
	var keys []string
	// 执行 SCAN 循环，每次查询 20 条符合 "temp*" 的键
	for {
		var err error
		keys, cursor, err = client.Scan(ctx, cursor, global.TempFile+"*", 20).Result()
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		// 处理返回的键名
		for _, key := range keys {
			result, err := client.HMGet(ctx, key, "expirationTime", "uploadDir").Result()
			if err != nil {
				break
			}
			fmt.Println(result[0])
			expirationTime := result[0].(string)
			// 将字符串转换为 int64
			intValue, _ := strconv.ParseInt(expirationTime, 10, 64)
			// 判断设定时间比现在时间小就删除
			if intValue < time.Now().Unix() {
				//删除图片
				fmt.Println("==================================")
				fmt.Println("删除了", result[1].(string))
				fmt.Println("==================================")
				b := imageUtils.DeleteImage(result[1].(string))
				fmt.Println(b)
				//删除redis图片消息
				client.Del(ctx, key)
			}
		}

		// 如果游标为 0 表示迭代完成
		if cursor == 0 {
			fmt.Println("每两分钟删除redis中的temp*的过期的键")
			time.Sleep(time.Minute * 2)
			*FileCh <- true
			break
		}
	}
	return
}
