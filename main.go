// 启动程序
package main

import (
	"context"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	v1 "golang-pet-api/api/v1"
	"golang-pet-api/global"
	"golang-pet-api/initialize"
	"golang-pet-api/routers"
	myValidator "golang-pet-api/validator"
	"os"
	"os/signal"
	"syscall"
)

// @title 宠物森林
// @version 1.0
// @description Go 语言编程之旅：一起用 Go 做项目
// @termsOfService https://github.com/zhangdonglin11/golang-pet-api
func main() {
	//1. 初始化logger
	initialize.InitLogger()
	//2. 初始化配置文件
	initialize.InitConfig()
	//3. 初始化routers
	router := routers.InitRouter()
	//4. 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}
	//5. 初始化数据库连接
	initialize.InitGorm()
	//6. 初始化redis连接
	initialize.InitRedis()

	// 启动websocket监听
	go v1.Manager.Start()

	// 初始化消息队列
	//initialize.InitRabbitMQ()
	// 处理队列信息 删除图片回显中的临时图片
	//go rabbitMQ.ReceiveImg()

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myValidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	// 定期删除临时文件
	go func() {
		ctx := context.Background()
		var FileCh chan bool
		FileCh = make(chan bool, 1)
		FileCh <- true
		for range FileCh {
			v1.HandleTimeOutFile(ctx, &FileCh)
		}
	}()

	/*
		1. S()可以获取一个全局的sugar，可以让我们自己设置一个全局的logger
		2. 日志是分级别的，debug， info ， warn， error， fetal
		3. S函数和L函数很有用， 提供了一个全局的安全访问logger的途径
	*/
	zap.S().Debugf("服务器启动，端口：%d", global.Config.Server.Address)
	err := router.Run(global.Config.Server.Address)
	if err != nil {
		zap.S().Panic("服务器启动，失败：", err.Error())
	}

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
