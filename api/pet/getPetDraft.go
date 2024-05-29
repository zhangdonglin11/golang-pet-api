package pet

import (
	"github.com/gin-gonic/gin"
	"golang-pet-api/models/forms"
	"golang-pet-api/service/pet_srv"
	"net/http"
)

// GetPetDraft 获取宠物草稿
func GetPetDraft(c *gin.Context) {
	uid, _ := c.Get("uid")
	userId := uid.(uint)
	// 查询宠物草稿
	petDraft, err := pet_srv.FindPetDraft(userId)
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":    0,
			"message": "没有宠物草稿",
		})
		return
	}

	data := forms.PetInfo{
		ID:            petDraft.ID,
		PetType:       petDraft.PetType,
		PetBreeds:     petDraft.PetBreeds,
		PetNickname:   petDraft.PetNickname,
		PetGender:     petDraft.PetGender,
		PetAge:        petDraft.PetAge,
		PetAddress:    petDraft.PetAddress,
		PetStatus:     petDraft.PetStatus,
		PetExperience: petDraft.PetExperience,
		PetAvatar:     petDraft.PetAvatar,
		PetIntro:      petDraft.PetIntro,
		Status:        petDraft.Status,
		UserID:        petDraft.UserID,
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    1,
		"message": "请求成功",
		"data":    data,
	})
	return
}
