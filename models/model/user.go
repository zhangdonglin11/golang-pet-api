package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Identity   string // 身份标识，可能是用户名或其他唯一标识
	Username   string `gorm:"unique"` // 用户名，唯一
	Password   string // 密码（这里应该存储加密后的密码哈希）
	Role       int    // 角色，0用户 1管理员
	Salt       string // 密码的盐
	ClientIp   string
	ClientPort string
	IsLogout   bool    // 是否已登出
	DeviceId   string  // 设备ID
	Profile    Profile // 用户信息（关联）
}

func (u User) GetTable() string {
	return "user"
}

type Profile struct {
	gorm.Model
	UserID   uint   `json:"userID"  gorm:"unique;not null"` // // 关联的用户ID，gorm的外键约束
	Nickname string `json:"nickname"`                       // 用户昵称
	Icon     string `json:"icon"`                           // 头像地址
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Gender   int    `json:"gender"` // 性别，1女 2男
	Status   int    `json:"status"` // 状态，例如是否有效
}

func (p Profile) GetTable() string {
	return "profile"
}
