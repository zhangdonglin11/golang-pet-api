package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PetId      uint       `json:"pet_id"`  // 宠物id
	UserID     uint       `json:"user_id"` // 用户id
	User       User       `json:"-"`
	TargetId   uint       `json:"target_id"` // 目标id
	TargetUser User       `gorm:"foreignKey:TargetId;references:ID"json:"-"`
	Level      int        `gorm:"size:4;default:0" json:"level"` // 层级 0 1
	RootID     uint       `json:"root_id"`                       // 根评论id
	SubComment *[]Comment `gorm:"foreignKey:RootID;references:ID" json:"-"`
	Content    string     `gorm:"type:text" json:"content"` // 评论内容
	Status     int        `gorm:"size:4" json:"status"`     // 字段状态
}
