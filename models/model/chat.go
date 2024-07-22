package model

import (
	"golang-pet-api/global"
	"gorm.io/gorm"
	"time"
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

// FindManyChatMessages 根据聊天ID和发送ID以及指定的起始时间，查询多条聊天消息。
func (ChatMessage) FindManyChatMessages(id, sendId string, startTime time.Time, pageSize int) ([]ChatMessage, error) {
	var chatMessages []ChatMessage
	if result := global.Db.Debug().
		Where("chat_id = ? OR chat_id = ?", id, sendId).
		Where(" created_at < ?", startTime).
		Order("created_at DESC").
		Limit(pageSize).
		Find(&chatMessages); result.Error != nil {
		return nil, result.Error
	}
	var chatMessagesList []ChatMessage
	for i := len(chatMessages) - 1; i >= 0; i-- {
		chatMessagesList = append(chatMessagesList, chatMessages[i])
	}
	return chatMessagesList, nil
}

// QueryUnreadMessagesBeforeTime 查询最新消息，传入pageSize=0侧查询endTime后的消息；传入pageSize不为0则查询最新的pageSize条
func (ChatMessage) QueryUnreadMessagesBeforeTime(id, sendId string, endTime time.Time, pageSize int) ([]ChatMessage, error) {
	var chatMessages []ChatMessage
	tx := global.Db.Debug().Model(ChatMessage{})
	tx.Where("chat_id = ? OR chat_id = ?", id, sendId)
	if pageSize == 0 {
		tx.Where("created_at > ?", endTime)
	} else {
		tx.Limit(pageSize)
	}
	if result := tx.Order("created_at DESC").
		Find(&chatMessages); result.Error != nil {
		return nil, result.Error
	}
	var chatMessagesList []ChatMessage
	for i := len(chatMessages) - 1; i >= 0; i-- {
		chatMessagesList = append(chatMessagesList, chatMessages[i])
	}
	return chatMessagesList, nil
}
