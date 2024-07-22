package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	v1 "golang-pet-api/api/v1"
	"golang-pet-api/docs"
	"golang-pet-api/global"
	"golang-pet-api/middleware"
	"net/http"
)

func InitRouter() *gin.Engine {
	// 设置模式
	gin.SetMode(global.Config.Server.Model)
	router := gin.Default()
	// 跌机恢复
	router.Use(gin.Recovery())
	router.Use(middleware.Cors())
	uploadDir := global.Config.ImageSettings.UploadDir
	router.StaticFS(uploadDir, http.Dir("./"+uploadDir))
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// v1
	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routerGroup := router.Group("/api/v1")

	//认证模块
	// 获取图片验证码
	routerGroup.GET("/captcha", v1.Auth{}.Captcha)
	// 获取短信验证码
	routerGroup.POST("/sms", v1.Auth{}.SmsCode)
	// 登录
	routerGroup.POST("/login", v1.Auth{}.Login)
	// 注册
	routerGroup.POST("/register", v1.Auth{}.Register)
	// 上传图片+回显
	routerGroup.POST("/upload", middleware.AuthMinddleware(), v1.UploadImageEcho)

	//用户模块
	userGroup := routerGroup.Group("/user", middleware.AuthMinddleware())
	// 获取用户信息
	userGroup.GET("/", v1.User{}.GetUser)
	// 修改用户信息
	userGroup.PUT("/", v1.User{}.UpdateUser)
	// 修改用户头像
	userGroup.PUT("/upload", v1.User{}.UploadPhoto)

	//宠物模块
	petGroup := routerGroup.Group("/pet")
	{
		//获取用户的宠物
		petGroup.GET("/myPet", middleware.AuthOptional(), v1.Pet{}.GetMyPet)
		// 取用户的宠物收藏列表
		petGroup.GET("/myLike", middleware.AuthOptional(), v1.Pet{}.GetMyLikePet)
		// 按条件获取所有宠物
		petGroup.POST("/filter", middleware.AuthOptional(), v1.Pet{}.GetListPet)
		// 获取一个宠物详细信息,无id则获取用户的草稿
		petGroup.GET("/:petId", middleware.AuthOptional(), v1.Pet{}.GetPet)
		// 创建和修改宠物信息
		petGroup.POST("/", middleware.AuthMinddleware(), v1.Pet{}.CreatePet)
		// 删除宠物
		petGroup.DELETE("/:petId", middleware.AuthMinddleware(), v1.Pet{}.DeletePet)
	}

	// 获取宠物的评论
	petGroup.GET("/comment", middleware.AuthOptional(), v1.Comment{}.GetCommentList)
	{
		// 获取评论的子评论
		petGroup.GET("/childComment", middleware.AuthOptional(), v1.Comment{}.GetCommentChild)
		// 提交宠物评论
		petGroup.POST("/submitComment", middleware.AuthMinddleware(), v1.Comment{}.CreateComment)
		// 删除宠物评论
		petGroup.DELETE("/comment/:cid", middleware.AuthMinddleware(), v1.Comment{}.DeleteComment)
	}

	// 一对一即时通讯模块
	chatGroup := routerGroup.Group("/chat")
	{
		// 创建用户关系列表
		chatGroup.GET("/:id", middleware.AuthMinddleware(), v1.BuildChatLIst)
		chatGroup.GET("/list", middleware.AuthMinddleware(), v1.GetChatList)
		chatGroup.GET("/wx", v1.WsHandler)
	}
	return router
}
