// 启动程序
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/global"
	"golang-pet-api/common/initialize"
	"golang-pet-api/common/utils/rabbitMQ"
	"golang-pet-api/routers"
)

func main() {
	// 初始化
	initialize.InitConfig()
	initialize.InitGorm()
	initialize.InitRedis()

	// 初始化消息队列
	initialize.InitRabbitMQ()
	// 处理队列信息 删除图片回显中的临时图片
	go rabbitMQ.ReceiveImg()

	// 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}
	// 创建数据库表

	// 设置模式
	gin.SetMode(global.Config.Server.Model)
	router := routers.InitRouter()

	addr := global.Config.Server.Address
	fmt.Println("服务器启动地址：", addr)
	err := router.Run(addr)
	if err != nil {
		fmt.Println("服务器启动失败！")
		return
	}
}
