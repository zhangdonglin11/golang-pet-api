package comment_srv

import (
	"errors"
	"golang-pet-api/common/global"
	"golang-pet-api/models/model"
	"golang-pet-api/service"
)

func GetPetCommentList(petId uint, level, page, size int) (total, pages int, comments []model.Comment) {
	// 查询评论总数
	var tempCount int64
	global.Db.Model(&model.Comment{PetId: petId}).Count(&tempCount)
	total = int(tempCount)
	// 评论有多少页
	if total == 0 {
		pages = 0
	} else {
		pages = total / size
		if pages < 1 {
			pages = 1
		}
	}

	global.Db.Scopes(service.Paginate(page, size)).
		Preload("User").
		Preload("TargetUser").
		Where(map[string]interface{}{"pet_id": petId, "level": level}).
		Order("created_at desc").
		Find(&comments)

	//global.Db.Scopes(service.Paginate(page, size)).Preload("User").
	//	Preload("TargetUser").
	//	Preload("SubComment", func(db *gorm.DB) *gorm.DB {
	//		return db.Limit(3).Preload("User").Preload("TargetUser").Order("created_at desc")
	//	}).
	//	Where(map[string]interface{}{"pet_id": petId, "level": level}).
	//	Order("created_at desc").
	//	Find(&comments)
	return
}

// 获取子评论
func GetPetSubCommentList() {

}

// 添加评论
func AddComment(comment *model.Comment) error {
	if result := global.Db.Create(comment); result.RowsAffected == 0 {
		return errors.New("添加评论失败")
	}
	return nil
}
