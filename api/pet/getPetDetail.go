package pet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/result"
	"golang-pet-api/models/forms"
	"golang-pet-api/service/pet_srv"
	"net/http"
	"strconv"
)

// GetOnePet 查询宠物或者宠物草稿 query pet_id 为0则查询草稿
func GetPetDetail(c *gin.Context) {
	uid, _ := c.Get("uid")
	userId := uid.(uint)
	//jwtUser, _ := tokenJWT.(*models.JwtUser)    // token对象
	petId, _ := strconv.Atoi(c.Param("pet_id")) //用户查询的宠物id
	fmt.Println("宠物id：", petId)
	// 查询目标宠物
	pet, err := pet_srv.FindPetById(uint(petId), 0)
	if err != nil {
		result.Failed(c, 501, err.Error())
		return
	}
	//获取宠物的点赞收藏信息
	count, isLike := pet_srv.GetLikeStatus(uint(petId), userId)

	data := map[string]interface{}{
		"pet": forms.PetInfo{
			ID:            pet.ID,
			PetType:       pet.PetType,
			PetBreeds:     pet.PetBreeds,
			PetNickname:   pet.PetNickname,
			PetGender:     pet.PetGender,
			PetAge:        pet.PetAge,
			PetAddress:    pet.PetAddress,
			PetStatus:     pet.PetStatus,
			PetExperience: pet.PetExperience,
			PetAvatar:     pet.PetAvatar,
			PetIntro:      pet.PetIntro,
			Status:        pet.Status,
			UserID:        pet.UserID,
		},
		"like": map[string]interface{}{
			"likeCount": count,
			"isLike":    isLike,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "请求成功",
		"data":    data,
	})
	return
}
