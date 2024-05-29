package pet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-pet-api/service/pet_srv"
	"net/http"
	"strconv"
)

func DeletePet(c *gin.Context) {
	uid, _ := c.Get("uid")
	userId := uid.(uint)
	petId, _ := strconv.Atoi(c.Query("pet_id"))
	fmt.Println(petId) //宠物id
	// 查询宠物
	_, err := pet_srv.FindPetById(uint(petId), userId)
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":    0,
			"message": "删除失败",
		})
		return
	}
	//删除宠物信息
	err = pet_srv.DeletePetByPid(uint(petId))
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":    0,
			"message": "删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    1,
		"message": "删除成功",
	})
	return
}
