package routers

import (
	"github.com/gin-gonic/gin"
	"golang-pet-api/api/auth"
	"golang-pet-api/api/comment"
	"golang-pet-api/api/pet"
	"golang-pet-api/api/user"
	"golang-pet-api/common/global"
	"golang-pet-api/middleware"
	"net/http"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	// 跌机恢复
	router.Use(gin.Recovery())
	router.Use(middleware.Cors())
	uploadDir := global.Config.ImageSettings.UploadDir
	router.StaticFS(uploadDir, http.Dir("./"+uploadDir))
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	register(router)
	return router
}

// 路由注册
func register(router *gin.Engine) {
	// 获取验证码
	router.GET("/api/captcha", auth.Captcha)
	// 登录
	router.POST("/api/login", auth.Login)
	// 注册
	router.POST("/api/register", auth.Register)

	userRouter := router.Group("/api/user", middleware.AuthMinddleware())
	// 获取用户信息
	userRouter.GET("/", user.GetUser)
	//// 修改用户信息
	userRouter.PUT("/", user.UpdateUser)
	userRouter.PUT("/upload", user.Upload)
	// 获取用户的宠物列表
	//userRouter.GET("/pet/:uid", GetMyPet)
	// 获取用户的宠物收藏列表
	//userRouter.GET("/like/:uid", GetLikePet)

	petRouter := router.Group("/api/pet")
	//petRouter.Use(middleware.AuthMinddleware())
	// 按条件获取所有宠物
	petRouter.POST("/filter", pet.GetPetList)
	// 获取宠物草稿
	petRouter.GET("/draft", middleware.AuthMinddleware(), pet.GetPetDraft)
	// 上传宠物图片回显
	petRouter.POST("/upload", middleware.AuthMinddleware(), pet.ImageEcho)
	// 获取一个宠物详细信息
	petRouter.GET("/:pet_id", middleware.AuthMinddleware(), pet.GetPetDetail)
	// 修改、保存宠物
	petRouter.PUT("/", middleware.AuthMinddleware(), pet.SavePet)
	// 删除宠物
	petRouter.DELETE("/", middleware.AuthMinddleware(), pet.DeletePet)

	// 获取宠物评论
	petRouter.GET("/comment", comment.GetPetComment)
	// 添加评论
	petRouter.POST("/comment", comment.AddComment)

}
