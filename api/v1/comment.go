package v1

import (
	"github.com/gin-gonic/gin"
	"golang-pet-api/global"
	"golang-pet-api/models/forms"
	"golang-pet-api/models/model"
	"golang-pet-api/utils"
	"strconv"
)

type Comment struct{}

// GetCommentList godoc
// @Summary 获取宠物的评论
// @Description 获取宠物的评论
// @Tags 评论模块
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param token header string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param pid query string true "宠物id"
// @Param page query string true "页码"
// @Param pageSize query string true "大小"
// @Success 200 {string} string "成功提交评论"
// @Failure 400 {string} string "无效的请求参数"
// @Router /api/v1/pet/comment [get]
func (m Comment) GetCommentList(c *gin.Context) {
	//f64 := c.GetFloat64("userId")
	//userId := uint(f64)

	pid, _ := c.GetQuery("pid")
	petId, _ := strconv.Atoi(pid)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	pages, comments, err := model.Comment{}.GetCommentByPetID(petId, page, pageSize)
	if err != nil {
		utils.RespFail(c.Writer, "内部错误："+err.Error())
		return
	}
	// 获取子评论
	for k, _ := range comments {
		temp := k
		_, childComment, _ := model.Comment{}.GetChildCommentByRootId(int(comments[temp].ID), 1, 3)
		comments[temp].ChildComments = &childComment
	}

	data := map[string]interface{}{
		"pages":    pages,
		"page":     page,
		"comments": comments,
	}

	utils.RespOk(c.Writer, data, "请求成功")
}

// GetCommentChild godoc
// @Summary 获取子评论
// @Description 获取宠物评论的子评论
// @Tags 评论模块
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param token header string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param cid query string true "评论id"
// @Param page query string true "页码"
// @Param pageSize query string true "大小"
// @Success 200 {string} string "成功提交评论"
// @Failure 400 {string} string "无效的请求参数"
// @Router /api/v1/pet/childComment [get]
func (m Comment) GetCommentChild(c *gin.Context) {
	cid, _ := c.GetQuery("cid")
	commentId, _ := strconv.Atoi(cid)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	pages, childComment, err := model.Comment{}.GetChildCommentByRootId(commentId, page, pageSize)
	if err != nil {
		utils.RespFail(c.Writer, "内部错误："+err.Error())
		return
	}
	data := map[string]interface{}{
		"childComment": childComment,
		"pages":        pages,
		"page":         page,
	}
	utils.RespOk(c.Writer, data, "请求成功")
}

// CreateComment godoc
// @Summary 提交评论
// @Description 提交新的评论
// @Tags 评论模块
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param token header string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param commentForm body forms.CommentForm true "评论表单"
// @Success 200 {string} string "成功提交评论"
// @Failure 400 {string} string "无效的请求参数"
// @Router /api/v1/pet/submitComment [post]
func (m Comment) CreateComment(c *gin.Context) {
	f64 := c.GetFloat64("userId")
	userId := uint(f64)

	var form forms.CommentForm
	if err := c.ShouldBind(&form); err != nil {
		HandleValidatorError(c, err)
		return
	}

	var comment model.Comment
	comment.PetId = form.PetId
	comment.TargetId = form.TargetId
	comment.Level = form.Level
	comment.RootId = form.RootID
	comment.Content = form.Content
	comment.UserId = userId
	comment.Status = 1

	if result := global.Db.Create(&comment); result.RowsAffected == 0 {
		utils.RespFail(c.Writer, "评论失败:"+result.Error.Error())
		return
	}

	utils.RespOk(c.Writer, form, "请求成功")
}

// DeleteComment godoc
// @Summary 删除评论
// @Description 删除评论
// @Tags 评论模块
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param token header string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param cid path string true "评论id"
// @Success 200 {string} string "成功提交评论"
// @Failure 400 {string} string "无效的请求参数"
// @Router /api/v1/pet/comment/{cid} [delete]
func (m Comment) DeleteComment(c *gin.Context) {
	f64 := c.GetFloat64("userId")
	userId := uint(f64)
	cid := c.Param("cid")

	if result := global.Db.Where("id = ? AND user_id", cid, userId).Delete(&model.Comment{}); result.RowsAffected == 0 {
		utils.RespFail(c.Writer, "删除评论失败")
		return
	}

	utils.RespOk(c.Writer, "", "请求成功")
}
