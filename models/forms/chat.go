package forms

// 发送消息的类型
type SendMsg struct {
	Content string `json:"content"`
	Type    int    `json:"type"`  // 1:发送信息 2：获取聊天记录 3：获取未读消息
	Media   int    `json:"media"` // 0: text, 1: image, 4: 表情
}

// 回复的消息
type ReplyMsg struct {
	Form    uint   `json:"form"`    // 发送者
	Code    int    `json:"code"`    // 状态码
	Msg     string `json:"msg"`     // 状态信息
	Type    int    `json:"type"`    // 1:发送信息 2：获取聊天记录 3：获取未读消息
	Media   int    `json:"media"`   // 0: text, 1: image, 4: 表情
	Content string `json:"content"` // 信息内容
}
