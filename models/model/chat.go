package model

import (
	"gorm.io/gorm"
)

type ChatList struct {
	gorm.Model
	ChatId         string `json:"chatId" gorm:"unique_index"` //聊天id
	UserId         uint   `json:"userId"`                     //用户id
	UserNickname   string `json:"userNickname"`
	UserIcon       string `json:"userIcon"`
	TargetId       uint   `json:"targetId"` //目标用户id
	TargetNickname string `json:"targetNickname"`
	TargetIcon     string `json:"targetIcon"`
	UnRead         int    `json:"unRead"`      //未读条数
	IsReceiving    bool   `json:"isReceiving"` //是否接收对方信息 拉黑功能 0false 1true
	Status         int    `json:"status"`      //字段状态 1显示列表，2隐藏列表
}

func (ChatList) TableName() string {
	return "chat_list"
}

type ChatMessage struct {
	gorm.Model
	ChatId  string `json:"chatId"`
	Media   int    `json:"media"` // 0: text, 1: image, 4: 表情
	Content string `json:"content"`
}

func (ChatMessage) TableName() string {
	return "chat_message"
}
