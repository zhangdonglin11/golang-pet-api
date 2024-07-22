package model

import (
	"golang-pet-api/global"
	"golang-pet-api/models"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	PetId  uint `json:"petId"`
	UserId uint `json:"userId"`
	//UserProfile   Profile `json:"userProfile" gorm:"foreignKey:UserId;references:UserID"`
	UserNickname string `json:"nickname"`
	UserIcon     string `json:"userIcon"`
	TargetId     uint   `json:"targetId"`
	//TargetProfile Profile `json:"targetProfile" gorm:"foreignKey:TargetId;references:UserID"`
	TargetNickname string `json:"targetNickname"`
	TargetIcon     string `json:"targetIcon"`
	Level          int    `json:"level"`
	RootId         uint   `json:"rootId"`
	//ChildComments *[]Comment `json:"childComments" gorm:"foreignKey:RootId;references:ID""`
	ChildComments *[]Comment `json:"childComments" gorm:"-"`
	Content       string     `json:"content"`
	Status        int        `json:"status"`
}

func (p Comment) GetTable() string {
	return "comment"
}

// GetCommentByPetID 获取根评论 宠物id,页面，大小
func (p Comment) GetCommentByPetID(petId, page, pageSize int) (pages int, comments []Comment, err error) {
	var count int64
	global.Db.Model(&Comment{}).Where("pet_id = ?", petId).Count(&count)
	pages = int(count) / page
	if pages < 1 {
		pages = 1
	}
	if result := global.Db.Debug().
		Table("comment").
		Select("comment.*, "+
			"user_profile.nickname as user_nickname, user_profile.icon as user_icon, "+
			"target_profile.nickname as target_nickname, target_profile.icon as target_icon").
		Joins("left join profile as user_profile on comment.user_id = user_profile.user_id").
		Joins("left join profile as target_profile on comment.target_id = target_profile.user_id").
		Where("comment.pet_id = ? AND comment.level = 0", petId).
		Scopes(models.Paginate(page, pageSize)).
		Find(&comments); result.Error != nil {
		return pages, nil, result.Error
	}
	return pages, comments, nil
}

// GetChildCommentByRootId 获取子评论 评论id,页面，大小
func (p Comment) GetChildCommentByRootId(cid, page, pageSize int) (pages int, childComment []Comment, err error) {
	var count int64
	global.Db.Model(&Comment{}).Where("root_id = ?", cid).Count(&count)
	pages = int(count) / page
	if pages < 1 {
		pages = 1
	}
	if result := global.Db.Debug().
		Table("comment").
		Select("comment.*, "+
			"user_profile.nickname as user_nickname, user_profile.icon as user_icon, "+
			"target_profile.nickname as target_nickname, target_profile.icon as target_icon").
		Joins("left join profile as user_profile on comment.user_id = user_profile.user_id").
		Joins("left join profile as target_profile on comment.target_id = target_profile.user_id").
		Where("comment.root_id = ?", cid).
		Scopes(models.Paginate(page, pageSize)).
		Find(&childComment); result.Error != nil {
		return pages, nil, result.Error
	}
	return pages, childComment, nil
}
