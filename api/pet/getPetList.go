package pet

import (
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/utils/jwt"
	"golang-pet-api/models/forms"
	"golang-pet-api/service/pet_srv"
	"net/http"
	"strings"
)

func GetPetList(c *gin.Context) {
	// 判断用户是否有token
	var userId uint
	authHeader := c.Request.Header.Get("token")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			// 验证token
			jwtUser, err := jwt.ValidateJWT(parts[1])
			if err == nil {
				userId = jwtUser.UserId
			}
		}
	}

	var petFilter forms.PetFilter
	c.ShouldBindJSON(&petFilter)

	var data []interface{}
	// 条件查询宠物信息
	petList := pet_srv.PetList(petFilter)
	for _, pet := range petList {
		// 查询该宠物的点赞信息
		count, isLike := pet_srv.GetLikeStatus(pet.ID, userId)
		// 处理每个宠物的详细信息
		data = append(data, map[string]interface{}{
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
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "请求成功",
		"data":    data,
	})
	return
}
