package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

// Pet 宠物信息模型
type Pet struct {
	gorm.Model
	PetType       int        `json:"petType" gorm:"size:4"`               //宠物类型
	PetBreeds     string     `json:"petBreeds" gorm:"size:10"`            //宠物瓶中
	PetNickname   string     `json:"petNickname" gorm:"size:10"`          //宠物昵称
	PetGender     string     `json:"petGender" gorm:"size:10"`            //宠物性别
	PetAge        string     `json:"petAge" gorm:"size:10"`               //宠物年龄
	PetAddress    string     `json:"petAddress" gorm:"size:64"`           //宠物所在地址
	PetStatus     string     `json:"petStatus" gorm:"size:10"`            //宠物状态
	PetExperience string     `json:"petExperience" gorm:"size:10"`        //宠物经验
	PetAvatar     AvatarList `json:"petAvatar" gorm:"type:varchar(1000)"` // 宠物图片字符串,空格分割
	PetIntro      string     `json:"petIntro" gorm:"type:text"`           //宠物简介
	Status        int        `json:"status" gorm:"size:4;not null"`       //字段状态 0草稿 1审核 2正常
	UserID        uint       `json:"userID" gorm:"not null"`              // 用户id
}

type AvatarList []string

// 实现 driver.Valuer 接口，Value 返回 json value
func (a AvatarList) Value() (driver.Value, error) {
	return json.Marshal(&a)
}

// 将数据库的json数据库转变为切片
func (a *AvatarList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &a)
}

// PetLike 宠物点赞模型
type PetLike struct {
	gorm.Model
	PetID  uint `json:"petID"`
	UserID uint `json:"userID"`
}
