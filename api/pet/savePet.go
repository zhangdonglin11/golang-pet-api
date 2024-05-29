package pet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/global"
	"golang-pet-api/common/utils/imageUtils"
	"golang-pet-api/models/forms"
	"golang-pet-api/models/model"
	"golang-pet-api/service/pet_srv"
	"net/http"
	"strings"
)

func SavePet(c *gin.Context) {
	uid, _ := c.Get("uid")
	userID := uid.(uint)
	// 提交的宠物数据
	var requestPet forms.PetInfo
	c.ShouldBindJSON(&requestPet)

	// 宠物信息发生变化只能这两种状态
	if requestPet.Status != 0 {
		requestPet.Status = 1
	}
	newPet := model.Pet{
		PetType:       requestPet.PetType,
		PetBreeds:     requestPet.PetBreeds,
		PetNickname:   requestPet.PetNickname,
		PetGender:     requestPet.PetGender,
		PetAge:        requestPet.PetAge,
		PetAddress:    requestPet.PetAddress,
		PetStatus:     requestPet.PetStatus,
		PetExperience: requestPet.PetExperience,
		PetAvatar:     requestPet.PetAvatar,
		PetIntro:      requestPet.PetIntro,
		Status:        requestPet.Status,
		UserID:        userID,
	}

	// 有宠物id 宠物id和用户id查询宠物是否存在
	if requestPet.ID != 0 {
		pet, err := pet_srv.FindPetById(requestPet.ID, userID)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code":    0,
				"message": "宠物保存失败！",
			})
			return
		}
		// 旧图片数组，更新宠物数据成功后删除
		oldAvatar := pet.PetAvatar
		// 新的图片文件名数组
		newAvatar := requestPet.PetAvatar

		// 更新宠物数据
		newPet.ID = requestPet.ID
		err = pet_srv.SavePet(&newPet)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "保存失败！",
			})
			return
		}

		// 将新的图片保存在一个map中，并删除redis保存的临时图片
		avatarMap := make(map[string]int)
		for _, v := range newAvatar {
			avatarMap[v] = 1
			// 输出redis保存的临时图片信息
			imageName := strings.TrimPrefix(v, global.Config.ImageSettings.ImageHost+global.Config.ImageSettings.UploadDir)
			s := imageUtils.TempRedisStore{}.Get(imageName, true)
			fmt.Println("redis-temp:", s)
		}

		// 遍历旧的图片数组 对比新的图片数组，删除旧图片数组的本地图片
		for _, v := range oldAvatar {
			if avatarMap[v] != 1 {
				// 获取图片文件名
				imageName := strings.TrimPrefix(v, global.Config.ImageSettings.ImageHost)
				if ok := imageUtils.DeleteImage(imageName); ok {
					fmt.Println("删除图片=", v)
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "保存成功！",
			"data": map[string]interface{}{
				"petId": newPet.ID,
			},
		})
		return
	}

	// 无宠物id 直接新建信息
	if requestPet.ID == 0 {
		err := pet_srv.CreatePet(&newPet)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "保存失败！",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "保存成功！",
				"data": map[string]interface{}{
					"petId": newPet.ID,
				},
			})
		}
		// 删除redis保存的临时图片
		for _, v := range newPet.PetAvatar {
			// 输出redis保存的临时图片信息
			imageName := strings.TrimPrefix(v, global.Config.ImageSettings.ImageHost+global.Config.ImageSettings.UploadDir)
			s := imageUtils.TempRedisStore{}.Get(imageName, true)
			fmt.Println("redis-temp:", s)
		}
		return
	}
}
