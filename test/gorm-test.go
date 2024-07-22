package main

import (
	"encoding/json"
	"fmt"
	"golang-pet-api/global"
	"golang-pet-api/initialize"
	"golang-pet-api/models/model"
)

func main() {
	initialize.InitConfig()
	initialize.InitGorm()
	//global.Db.AutoMigrate(&model.Comment{})
	//var comments []model.Comment
	//if result := global.Db.Debug().Model(&model.Comment{}).
	//	Where("pet_id = ? AND level=0", 1).
	//	Preload("UserProfile").
	//	Preload("TargetProfile").
	//	Preload("Comments").
	//	Preload("Comments.UserProfile").
	//	Preload("Comments.TargetProfile").
	//	Scopes(models.Paginate(1, 5)).
	//	Find(&comments); result.Error != nil {
	//	fmt.Println(result.Error)
	//}

	//if result := global.Db.Debug().Model(&model.Comment{}).
	//	Where("pet_id = ? AND level=0", 1).
	//	Preload("UserProfile").
	//	Preload("TargetProfile").
	//	Scopes(models.Paginate(1, 5)).
	//	Find(&comments); result.Error != nil {
	//	fmt.Println(result.Error)
	//}
	//
	//for k, _ := range comments {
	//	var childComment []model.Comment
	//	temp := k
	//	if result := global.Db.Debug().Model(&model.Comment{}).
	//		Where("root_id", comments[temp].ID).
	//		Preload("UserProfile").
	//		Preload("TargetProfile").
	//		Scopes(models.Paginate(1, 5)).
	//		Find(&childComment); result.Error != nil {
	//		fmt.Println(result.Error)
	//	}
	//	comments[temp].Comments = &childComment
	//}

	//var comments []model.Comment

	//if result := global.Db.Debug().
	//	Table("comment").
	//	Select("comment.*, "+
	//		"user_profile.nickname as user_nickname, user_profile.icon as user_icon, "+
	//		"target_profile.nickname as target_nickname, target_profile.icon as target_icon").
	//	Joins("left join profile as user_profile on comment.user_id = user_profile.user_id").
	//	Joins("left join profile as target_profile on comment.target_id = target_profile.user_id").
	//	Where("comment.pet_id = ? AND comment.level = 0", 1).
	//	Scopes(models.Paginate(1, 5)).
	//	Find(&comments); result.Error != nil {
	//	fmt.Println(result.Error)
	//}
	//for k, _ := range comments {
	//	temp := k
	//	var childComment []model.Comment
	//	if result := global.Db.Debug().
	//		Table("comment").
	//		Select("comment.*, "+
	//			"user_profile.nickname as user_nickname,user_profile.icon as user_icon,"+
	//			"target_profile.nickname as target_nickname,target_profile.icon as target_icon").
	//		Joins("left join profile as user_profile on comment.user_id = user_profile.user_id").
	//		Joins("left join profile as target_profile on comment.target_id = target_profile.user_id").
	//		Where("root_id", comments[temp].ID).
	//		Scopes(models.Paginate(1, 3)).
	//		Find(&childComment); result.Error != nil {
	//		fmt.Println(result.Error)
	//	}
	//	comments[temp].ChildComments = &childComment
	//}

	//comments, err := model.Comment{}.GetCommentByPetID(1, 1, 5)
	//if err != nil {
	//	return
	//}
	//for k, _ := range comments {
	//	temp := k
	//	childComment, err := model.Comment{}.GetChildCommentByRootId(int(comments[temp].ID), 1, 3)
	//	if err != nil {
	//		return
	//	}
	//	comments[temp].ChildComments = &childComment
	//}
	//
	//marshal, err := json.Marshal(comments)
	//if err != nil {
	//	return
	//}

	var chatLists []model.ChatList
	if err := global.Db.Table("chat_list").
		Select("chat_list.*, user_profile.nickname AS user_nickname, user_profile.icon AS user_icon, target_profile.nickname AS target_nickname, target_profile.icon AS target_icon").
		Joins("LEFT JOIN profile AS user_profile ON user_profile.user_id = chat_list.user_id").
		Joins("LEFT JOIN profile AS target_profile ON target_profile.user_id = chat_list.target_id").
		Where("chat_list.user_id = ?", 3). // 替换为你的用户ID
		Find(&chatLists).Error; err != nil {

	}
	marshal, err := json.Marshal(&chatLists)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}
