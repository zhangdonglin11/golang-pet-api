package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-pet-api/global"
	"golang-pet-api/models"
	"golang-pet-api/models/forms"
	"golang-pet-api/models/model"
	"golang-pet-api/utils"
	"golang-pet-api/utils/imageUtils"
	"strconv"
)

type Pet struct{}

// GetPet godoc
// @Summary 获取宠物详细信息
// @Description 根据宠物id获取宠物信息：[get] /api/v1/pet/:petId
// @Tags 宠物模块
// @Produce json
// @Security ApiKeyAuth
// @Param token header string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param petId path string false "宠物id"
// @Success 200 {string} json{Code,Msg,Data}  "成功"
// @Router /api/v1/pet/{petId} [get]
func (p Pet) GetPet(c *gin.Context) {
	petId := c.Param("petId")
	f64 := c.GetFloat64("userId")
	userId := uint(f64)

	var pet model.Pet
	if result := global.Db.First(&pet, petId); result.RowsAffected == 0 {
		utils.RespFail(c.Writer, "宠物不存在")
		return
	}

	var count int64
	var isLike bool
	var petLike model.PetLike
	global.Db.Find(&petLike, petId).Count(&count)
	if result := global.Db.Where("pet_id = ? AND user_id", petId, userId).First(&petLike); result.RowsAffected != 0 {
		isLike = true
	}
	petLikeMap := map[string]interface{}{
		"count":  count,
		"isLike": isLike,
	}

	data := map[string]interface{}{
		"petInfo": pet,
		"petLike": petLikeMap,
	}
	utils.RespOk(c.Writer, data, "请求成功")
}

// GetListPet godoc
// @Summary 通过条件获取宠物列表
// @Description 描述
// @Tags 宠物模块
// @Produce json
// @Security ApiKeyAuth
// @Param token      header   string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param petType    formData int    false "1:狗|2:猫"
// @Param petAddress formData string false "宠物地址"
// @Param petStatus  formData string false "宠物社交状态"
// @Param petGender  formData string false "性别"
// @Param petAge     formData string false "年龄"
// @Param page 		 formData    string true  "查询第几页"
// @Param pageSize   formData    string true  "每页多少条数据"
// @Success 200 {string} json{Code,Msg,Data}  "成功"
// @Router /api/v1/pet/filter [post]
func (p Pet) GetListPet(c *gin.Context) {
	// 获取token解析当前用户的id，token可能没有
	f64 := c.GetFloat64("userId")
	userId := uint(f64)

	// 绑定请求的数据
	var form forms.FilterForm
	if err := c.ShouldBind(&form); err != nil {
		HandleValidatorError(c, err)
		return
	}

	tx := global.Db.Model(model.Pet{})

	if form.PetType != 0 {
		tx.Where("pet_type = ?", form.PetType)
	}
	if form.PetAddress != "" {
		tx.Where("pet_address LIKE %?%", form.PetAddress)
	}
	if form.PetStatus != "" {
		tx.Where("pet_status = ?", form.PetStatus)
	}
	if form.PetGender != "" {
		tx.Where("pet_gender = ?", form.PetGender)
	}
	if form.PetAge != "" {
		tx.Where("pet_age = ?", form.PetAge)
	}
	var pets []model.Pet
	//总页数
	var count int64
	tx.Count(&count)
	pages := int(count) / form.Page
	if pages < 1 {
		pages = 1
	}
	// 分页查询
	tx.Scopes(models.Paginate(form.Page, form.PageSize)).Find(&pets)

	var petArr []interface{}
	for _, pet := range pets {
		petMap := make(map[string]interface{})

		// 查询宠物点赞数量和状态
		var count int64
		var isLike bool
		var petLike model.PetLike
		global.Db.Model(model.PetLike{}).Where("pet_id = ?", pet.ID).Count(&count)
		//有用户id才能查询到结果
		if result := global.Db.Model(model.PetLike{}).Where("pet_id=? AND user_id = ?", pet.ID, userId).First(&petLike); result.RowsAffected != 0 {
			isLike = true
		}

		petLikeMap := map[string]interface{}{
			"count":  count,
			"isLike": isLike,
		}

		petMap["petLike"] = petLikeMap
		petMap["petInfo"] = pet
		petArr = append(petArr, petMap)
	}

	data := map[string]interface{}{
		"petList": petArr,
		"pages":   pages,
		"page":    form.Page,
	}

	utils.RespOk(c.Writer, data, "请求成功")
	return
}

// GetMyPet godoc
// @Summary 获取我的宠物
// @Description 查询目标收藏的宠物列表；[get] /api/v1/pet/myPet?userId=,page=,pageSize=
// @Tags 宠物模块
// @Produce json
// @Security ApiKeyAuth
// @Param token header string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param userId query string true "目标用户id"
// @Param page query string true "查询第几页"
// @Param pageSize query string true "每页多少条数据"
// @Success 200 {string} json{Code,Msg,Data}  "成功"
// @Router /api/v1/pet/myPet [get]
func (p Pet) GetMyPet(c *gin.Context) {
	f64 := c.GetFloat64("userId")
	userId := uint(f64)

	// 查询目标id收藏的宠物列表
	targetID, _ := c.GetQuery("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	//总页数
	var count int64
	global.Db.Model(&model.Pet{}).Where("user_id=?", targetID).Count(&count)
	pages := int(count) / page
	if pages < 1 {
		pages = 1
	}
	//根据用户id批量查询宠物
	var pets []model.Pet
	global.Db.Where("user_id=?", targetID).Scopes(models.Paginate(page, pageSize)).Find(&pets)

	//处理返回的信息
	var petArr []interface{}
	for _, pet := range pets {
		petMap := make(map[string]interface{})
		// 查询宠物点赞数量和状态
		var count int64
		var isLike bool
		var petLike model.PetLike
		global.Db.Model(model.PetLike{}).Where("pet_id = ?", pet.ID).Count(&count)
		//有用户id才能查询到结果
		if userId != 0 {
			if result := global.Db.Model(model.PetLike{}).Where("pet_id=? AND user_id = ?", pet.ID, userId).First(&petLike); result.RowsAffected != 0 {
				isLike = true
			}
		}
		petLikeMap := map[string]interface{}{
			"count":  count,
			"isLike": isLike,
		}
		petMap["petLike"] = petLikeMap
		petMap["petInfo"] = pet
		petArr = append(petArr, petMap)
	}
	data := map[string]interface{}{
		"petList": petArr,
		"pages":   pages,
		"page":    page,
	}
	utils.RespOk(c.Writer, data, "请求成功")
	return
}

// GetMyLikePet godoc
// @Summary 获取我的宠物收藏
// @Description 查询目标收藏的宠物列表；[get] /api/v1/pet/myLike?userId=,page=,pageSize=
// @Tags 宠物模块
// @Produce json
// @Security ApiKeyAuth
// @Param token header string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param userId query string true "目标用户id"
// @Param page query string true "查询第几页"
// @Param pageSize query string true "每页多少条数据"
// @Success 200 {string} json{Code,Msg,Data}  "成功"
// @Router /api/v1/pet/myLike [get]
func (p Pet) GetMyLikePet(c *gin.Context) {
	// 当前用户的id
	f64 := c.GetFloat64("userId")
	userId := uint(f64)
	// 查询目标id收藏的宠物列表
	targetID, _ := c.GetQuery("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 分页总数
	var count int64
	global.Db.Model(&model.PetLike{}).Where("user_id = ?", targetID).Count(&count)
	pages := int(count) / page
	if pages < 1 {
		pages = 1
	}
	// 获取目标id 的点赞的宠物id切片
	var petIdArr []uint
	var petLikeList []model.PetLike
	global.Db.Where("user_id = ?", targetID).Scopes(models.Paginate(page, pageSize)).Find(&petLikeList)
	for _, v := range petLikeList {
		petIdArr = append(petIdArr, v.PetID)
	}

	// 根据宠物id切片批量查询宠物
	var pets []model.Pet
	global.Db.Where("id in (?)", petIdArr).Find(&pets)

	var petArr []interface{}
	for _, pet := range pets {
		petMap := make(map[string]interface{})

		// 查询宠物点赞数量和状态
		var isLike bool
		var petLike model.PetLike
		global.Db.Model(model.PetLike{}).Where("pet_id = ?", pet.ID).Count(&count)
		// 根据宠物id和用户id查询收藏状态，用户id为0则查询不到
		if result := global.Db.Model(model.PetLike{}).Where("pet_id=? AND user_id = ?", pet.ID, userId).First(&petLike); result.RowsAffected != 0 {
			isLike = true
		}
		petLikeMap := map[string]interface{}{
			"count":  count,
			"isLike": isLike,
		}

		petMap["petLike"] = petLikeMap
		petMap["petInfo"] = pet
		petArr = append(petArr, petMap)
	}

	data := map[string]interface{}{
		"petList": petArr,
		"pages":   pages,
		"page":    page,
	}
	utils.RespOk(c.Writer, data, "请求成功")
	return
}

// CreatePet godoc
// @Summary 创建或修改宠物
// @Description 创建或修改宠物信息，宠物id为空则创建，不为空则修改。
// @Tags 宠物模块
// @Produce json
// @Security ApiKeyAuth
// @Param token header string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param petId formData int false "宠物ID (为空则创建新宠物)"
// @Param petType formData int true "宠物类型 (1:狗, 2:猫)" enum(1, 2)
// @Param petBreeds formData string false "宠物品种"
// @Param petNickname formData string false "宠物昵称"
// @Param petGender formData string false "宠物性别"
// @Param petAge formData string false "宠物年龄"
// @Param petAddress formData string false "宠物地址"
// @Param petStatus formData string false "宠物状态"
// @Param petExperience formData string false "宠物经历"
// @Param petAvatar formData []string false "宠物头像链接数组"
// @Param petIntro formData string false "宠物介绍"
// @Param status formData int true "状态 (1:草稿, 2:发布)" enum(1, 2)
// @Success 200 {string} json {"Code":0, "Msg":"成功", "Data":{}}
// @Router /api/v1/pet [post]
func (p Pet) CreatePet(c *gin.Context) {
	f64 := c.GetFloat64("userId")
	userId := uint(f64)

	var form forms.PetForm
	if err := c.ShouldBind(&form); err != nil {
		HandleValidatorError(c, err)
		return
	}
	// 修改宠物
	if form.PetId != 0 {
		var pet model.Pet
		if result := global.Db.Where("id=? AND user_id=?", form.PetId, userId).First(&pet); result.RowsAffected == 0 {
			utils.RespFail(c.Writer, "宠物不存在")
			return
		}
		//保存就图片的信息
		orderPetAvatar := pet.PetAvatar

		pet.PetType = form.PetType
		pet.PetBreeds = form.PetBreeds
		pet.PetNickname = form.PetNickname
		pet.PetGender = form.PetGender
		pet.PetAge = form.PetAge
		pet.PetAddress = form.PetAddress
		pet.PetStatus = form.PetStatus
		pet.PetExperience = form.PetExperience
		pet.PetIntro = form.PetIntro
		pet.Status = form.Status
		pet.PetAvatar = form.PetAvatar

		if result := global.Db.Save(&pet); result.RowsAffected == 0 {
			utils.RespFail(c.Writer, "宠物保存失败")
			return
		}

		// 图片新旧信息处理
		petAvatarArr := append(orderPetAvatar, form.PetAvatar...)
		//并集
		union := make(map[string]int)
		// 交集
		intersection := make(map[string]int)
		for _, key := range petAvatarArr {
			if union[key] == 1 {
				intersection[key] = 1
			} else {
				union[key] = 1
			}
		}
		// 遍历旧图片切片,判断是否在交集overlapMap[key]中,不在就删除的旧图片
		for _, val := range orderPetAvatar {
			if intersection[val] != 1 {
				// 获取文件保存路径
				path, err := imageUtils.GetUrlByFilePath(val)
				if err != nil {
					fmt.Println(err)
					continue
				}
				imageUtils.DeleteImage(path)
			}
		}
		// 遍历新图片切片,判断是否在交集overlapMap[key]中,不在就删除redis中的图片信息
		for _, url := range form.PetAvatar {
			if intersection[url] != 1 {
				//获取url中的文件名
				fileName, err := imageUtils.GetUrlByFileName(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				//删除redis中的信息
				global.RedisDb.Del(c, global.TempFile+fileName)
			}
		}

		utils.RespOk(c.Writer, pet, "宠物保存成功")
		return
	}

	//无宠物id 创建
	var pet model.Pet
	pet.PetType = form.PetType
	pet.PetBreeds = form.PetBreeds
	pet.PetNickname = form.PetNickname
	pet.PetGender = form.PetGender
	pet.PetAge = form.PetAge
	pet.PetAddress = form.PetAddress
	pet.PetStatus = form.PetStatus
	pet.PetExperience = form.PetExperience
	pet.PetAvatar = form.PetAvatar
	pet.PetIntro = form.PetIntro
	pet.Status = form.Status
	//pet.UserID = userId

	if result := global.Db.Create(&pet); result.RowsAffected == 0 {
		utils.RespFail(c.Writer, "宠物保存失败")
		return
	}

	// 保存宠物成功后,删除redis的图片信息
	for _, avatar := range form.PetAvatar {
		//获取url中的文件名
		fileName, err := imageUtils.GetUrlByFileName(avatar)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//删除redis中的信息
		global.RedisDb.Del(c, global.TempFile+fileName)
	}
	utils.RespOk(c.Writer, form, "请求成功")
	return
}

// DeletePet godoc
// @Summary 删除宠物信息
// @Description  根据宠物id删除宠物信息：[delete] /api/v1/pet/:petId
// @Tags 宠物模块
// @Produce json
// @Security ApiKeyAuth
// @Param token header string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param petId path string true "宠物id"
// @Success 200 {string} json{Code,Msg,Data} "成功"
// @Router /api/v1/pet/{petId} [delete]
func (p Pet) DeletePet(c *gin.Context) {
	f64 := c.GetFloat64("userId")
	userId := uint(f64)
	petId := c.Param("petId")

	tx := global.Db.Begin()
	var pet model.Pet
	if result := tx.Where("id=? AND user_id=?", petId, userId).First(&pet); result.RowsAffected == 0 {
		tx.Callback()
		utils.RespFail(c.Writer, "未找到符合条件的宠物进行删除")
		return
	}
	if result := tx.Delete(&model.Pet{}, pet.ID); result.Error != nil {
		tx.Callback()
		utils.RespFail(c.Writer, result.Error.Error())
		return
	}
	//删除收藏表
	if result := tx.Where("pet_id=?", pet.ID).Delete(&model.PetLike{}); result.Error != nil {
		tx.Callback()
		utils.RespFail(c.Writer, result.Error.Error())
		return
	}
	//删除评论表
	if result := tx.Where("pet_id=?", pet.ID).Delete(&model.Comment{}); result.Error != nil {
		tx.Callback()
		utils.RespFail(c.Writer, result.Error.Error())
		return
	}

	tx.Commit()
	// 删除图片
	for _, val := range pet.PetAvatar {
		path, err := imageUtils.GetUrlByFilePath(val)
		if err != nil {
			continue
		}
		imageUtils.DeleteImage(path)
	}
	utils.RespOk(c.Writer, "", "请求成功")
	return
}
