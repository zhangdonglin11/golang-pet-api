package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang-pet-api/common/global"
	"golang-pet-api/models/model"
	"golang-pet-api/service/comment_srv"
	"golang-pet-api/service/pet_srv"
	"net/http"
	"strings"
)

type FormData struct {
	PetId    uint   `json:"pet_id" binding:"required"`
	TargetId uint   `json:"target_id"`
	Level    int    `json:"level"`
	RootID   uint   `json:"root_id"`
	Content  string `json:"content" binding:"required"`
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// AddComment 添加评论
func AddComment(c *gin.Context) {
	uid, _ := c.Get("uid")
	userId := uid.(uint)

	var formData FormData
	if err := c.ShouldBind(&formData); err != nil {
		HandleValidatorError(c, err)
		return
	}
	if _, err := pet_srv.FindPetById(formData.PetId, 0); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": err,
		})
		return
	}

	comment := model.Comment{
		TargetId: formData.TargetId,
		Level:    formData.Level,
		RootID:   formData.RootID,
		Content:  formData.Content,
		UserID:   userId,
	}
	if err := comment_srv.AddComment(&comment); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "添加评论成功",
	})
	return

}
