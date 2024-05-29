package comment

import (
	"github.com/gin-gonic/gin"
	"golang-pet-api/common/global"
	"golang-pet-api/models/model"
	"golang-pet-api/service/comment_srv"
	"net/http"
)

// 利用递归 整理评论数据
func handleComment(comments *[]model.Comment) []map[string]interface{} {
	var data []map[string]interface{}
	for _, comment := range *comments {
		var subComment []map[string]interface{}
		if comment.SubComment != nil {
			subComment = handleComment(comment.SubComment)
		}
		var count int64
		global.Db.Model(&model.Comment{}).
			Where(map[string]interface{}{"pet_id": comment.PetId, "root_id": comment.RootID}).
			Scan(&count)
		commentMap := map[string]interface{}{
			"comment_id":       comment.ID,
			"pet_id":           comment.PetId,
			"user_id":          comment.User.ID,
			"user_name":        comment.User.NickName,
			"target_id":        comment.TargetId,
			"target_name":      comment.TargetUser.NickName,
			"level":            comment.Level,
			"root_id":          comment.RootID,
			"sub_comment":      subComment,
			"sub_commentTotal": len(subComment),
			"content":          comment.Content,
			"status":           comment.Status,
		}
		data = append(data, commentMap)
	}
	return data
}

type Condition struct {
	Page  int  `json:"page"`
	Size  int  `json:"size"`
	PetId uint `json:"petId"`
}

// GetPetComment 获取评论
func GetPetComment(c *gin.Context) {
	var condition Condition
	c.ShouldBindJSON(&condition)

	total, pages, comments := comment_srv.GetPetCommentList(condition.PetId, 0, condition.Page, condition.Size)

	data := map[string]interface{}{
		"total":   total,
		"pages":   pages,
		"page":    condition.Page,
		"comment": handleComment(&comments),
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "请求成功",
		"data":    data,
	})
	return
}
