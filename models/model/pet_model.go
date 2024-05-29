package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

// Pet 宠物信息模型
type Pet struct {
	gorm.Model
	PetType       int        `gorm:"size:4"`             //宠物类型
	PetBreeds     string     `gorm:"size:10"`            //宠物瓶中
	PetNickname   string     `gorm:"size:10"`            //宠物昵称
	PetGender     string     `gorm:"size:10"`            //宠物性别
	PetAge        string     `gorm:"size:10"`            //宠物年龄
	PetAddress    string     `gorm:"size:64"`            //宠物所在地址
	PetStatus     string     `gorm:"size:10"`            //宠物状态
	PetExperience string     `gorm:"size:10"`            //宠物经验
	PetAvatar     AvatarList `gorm:"type:varchar(1000)"` // 宠物图片字符串,空格分割
	PetIntro      string     `gorm:"type:text"`          //宠物简介
	Status        int        `gorm:"size:4;not null"`    //字段状态 0草稿 1审核 2正常
	UserID        uint       `gorm:"not null"`           // 用户id
}

type AvatarList []string

// 实现 driver.Valuer 接口，Value 返回 json value
func (a AvatarList) Value() (driver.Value, error) {
	return json.Marshal(&a)
}

// 将数据库的json数据库转变为切边
func (a *AvatarList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &a)
}

// PetLike 宠物点赞模型
type PetLike struct {
	gorm.Model
	PetID  uint
	UserID uint
}
