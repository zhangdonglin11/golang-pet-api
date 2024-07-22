package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang-pet-api/global"
	"golang-pet-api/models/forms"
	"golang-pet-api/models/model"
	"golang-pet-api/utils"
	"golang-pet-api/utils/e"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const month = 60 * 60 * 24 * 30 // 按照30天算一个月

// BuildChatLIst godoc
// @Summary 创建用户聊天关系
// @Description 根据目标用户id建立聊天关系：[get] /api/v1/chat/:id
// @Tags 聊天模块
// @Produce json
// @Security ApiKeyAuth
// @Param token header string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param id path string false "目标用户id"
// @Success 200 {string} json{Code,Msg,Data}  "成功"
// @Router /api/v1/chat/{id} [get]
func BuildChatLIst(ctx *gin.Context) {
	useridF64 := ctx.GetFloat64("userId")

	targetedString := ctx.Param("id")
	targetedInt, _ := strconv.Atoi(targetedString)

	tx := global.Db.Begin()
	//查询用户列表
	userId := createId(useridF64, targetedString)
	var chatList1 model.ChatList
	if result := global.Db.Where(model.ChatList{ChatId: userId}).First(&chatList1); result.RowsAffected == 0 {
		chatList1.ChatId = userId
		chatList1.UserId = uint(useridF64)
		chatList1.TargetId = uint(targetedInt)
		chatList1.IsReceiving = true
		chatList1.Status = 1
		if tx.Create(&chatList1).RowsAffected == 0 {
			tx.Callback()
			utils.RespFail(ctx.Writer, "请求失败")
			return
		}
	}
	targetId := createId(targetedString, useridF64)
	var chatList2 model.ChatList
	if result := global.Db.Where(model.ChatList{ChatId: targetId}).First(&chatList2); result.RowsAffected == 0 {
		chatList2.ChatId = targetId
		chatList2.UserId = uint(targetedInt)
		chatList2.TargetId = uint(useridF64)
		chatList2.IsReceiving = true
		chatList2.Status = 1
		if tx.Create(&chatList2).RowsAffected == 0 {
			tx.Callback()
			utils.RespFail(ctx.Writer, "请求失败")
			return
		}
	}
	//判断用户是否逻辑删除过聊天列表
	if chatList1.Status == 0 {
		chatList1.Status = 1
		tx.Save(&chatList1)
	}
	tx.Commit()
	utils.RespOk(ctx.Writer, nil, "请求成功")
	return
}

// GetChatList godoc
// @Summary 获取用户聊天列表
// @Description 根据目标用户id建立聊天关系：[get] /api/v1/chat/:id
// @Tags 聊天模块
// @Produce json
// @Security ApiKeyAuth
// @Param token header string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Success 200 {string} json{Code,Msg,Data}  "成功"
// @Router /api/v1/chat/list [get]
func GetChatList(ctx *gin.Context) {
	userId_f64 := ctx.GetFloat64("userId")
	var chatLists []model.ChatList
	if err := global.Db.Table("chat_list").
		Select("chat_list.*, user_profile.nickname AS user_nickname, user_profile.icon AS user_icon, target_profile.nickname AS target_nickname, target_profile.icon AS target_icon").
		Joins("LEFT JOIN profile AS user_profile ON user_profile.user_id = chat_list.user_id").
		Joins("LEFT JOIN profile AS target_profile ON target_profile.user_id = chat_list.target_id").
		Where("chat_list.user_id = ?", userId_f64). // 替换为你的用户ID
		Find(&chatLists).Error; err != nil {
		utils.RespFail(ctx.Writer, "请求错误")
		return
	}
	utils.RespOk(ctx.Writer, chatLists, "请求成功")
	return
}

// WsHandler godoc
// @Summary 建立websocket连接
// @Description 使用postman或apifox等接口软件使用websocket进行接口测试，根据目标用户id建立聊天长连接：[ws] ws://127.0.0.1:8088/api/v1/chat/wx?uid=2
// @Tags 聊天模块
// @Produce json
// @Security ApiKeyAuth
// @Param token header string false "Bearer Token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3NTQxNzksImlhdCI6MTcyMTAyNjE3OSwicm9sZSI6MCwidXNlcklkIjoxfQ.7_QhxomXnG1TLwggvkij1UwJkaCxxFtldUEvzbWbHWM)
// @Param id path string false "目标用户id"
// @Success 200 {string} json{Code,Msg,Data}  "成功"
// @Router /api/v1/chat/wx [get]
func WsHandler(c *gin.Context) {
	userId_f64 := c.GetFloat64("userId")
	uid_string := strconv.Itoa(int(userId_f64))

	targetId_string := c.Query("toUid")

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
			return true
		}}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	id := createId(uid_string, targetId_string)
	sendID := createId(targetId_string, uid_string)
	// 创建一个用户实例
	client := &Client{
		ID:     id,
		SendID: sendID,
		Socket: conn,
		Send:   make(chan []byte),
	}
	// 用户注册到用户管理上
	Manager.Register <- client
	go client.Read(c)
	go client.Write()
}

func createId(uid, toUid interface{}) string {
	sprintf := fmt.Sprintf("%v->%v", uid, toUid)
	return sprintf
}
func resolveId(sendId string) (senderId, recipientId uint, err error) {
	parts := strings.Split(sendId, "->")
	if len(parts) != 2 {
		return 0, 0, errors.New("sendId format error")
	}
	uid, _ := strconv.Atoi(parts[0])
	toUid, _ := strconv.Atoi(parts[1])
	return uint(uid), uint(toUid), nil
}

func (c *Client) Read(ctx *gin.Context) {
	defer func() { //避免忘记关闭，使用defer保险关闭
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		sendMsg := new(forms.SendMsg)
		err := c.Socket.ReadJSON(&sendMsg)

		// 读取json格式，如果不是json格式，会报错
		if err != nil {
			log.Println("数据格式不正确", err)
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}

		if sendMsg.Type == 1 { // 发送信息
			r1, _ := global.RedisDb.Get(ctx, c.ID).Result()
			r2, _ := global.RedisDb.Get(ctx, c.SendID).Result()
			senderId, _, _ := resolveId(c.ID)
			if r1 >= "3" && r2 == "0" {
				replyMsg := forms.ReplyMsg{
					Code: e.WebsocketLimit,
					Msg:  "到达限制",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				_, _ = global.RedisDb.Expire(ctx, c.ID, time.Hour*24*30).Result()
				continue
			} else {
				global.RedisDb.Incr(ctx, c.ID)
				_, _ = global.RedisDb.Expire(ctx, c.ID, time.Hour*24*30).Result()
			}
			//处理返回信息体
			replyMsg := forms.ReplyMsg{
				Form:    senderId,
				Code:    e.WebsocketSuccessMessage,
				Content: sendMsg.Content,
				Media:   sendMsg.Media,
				Type:    sendMsg.Type,
			}
			log.Println(c.ID, "：发送信息", replyMsg)
			msg, _ := json.Marshal(&replyMsg)
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: msg,
			}
		} else if sendMsg.Type == 2 { //拉取历史消息
			//timeT, err := strconv.Atoi(sendMsg.Content) // 传送来时间
			//if err != nil {
			//	timeT = 999999999
			//}
			//results, _ := FindMany(conf.MongoDBName, c.SendID, c.ID, int64(timeT), 10)
			//if len(results) > 10 {
			//	results = results[:10]
			//} else if len(results) == 0 {
			//	replyMsg := ReplyMsg{
			//		Code:    e.WebsocketEnd,
			//		Content: "到底了",
			//	}
			//	msg, _ := json.Marshal(replyMsg)
			//	_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			//	continue
			//}
			//for _, result := range results {
			//	replyMsg := ReplyMsg{
			//		From:    result.From,
			//		Content: fmt.Sprintf("%s", result.Msg),
			//	}
			//	msg, _ := json.Marshal(replyMsg)
			//	_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			//}
		} else if sendMsg.Type == 3 {
			//results, err := FirsFindtMsg(conf.MongoDBName, c.SendID, c.ID)
			//if err != nil {
			//	log.Println(err)
			//}
			//for _, result := range results {
			//	replyMsg := ReplyMsg{
			//		From:    result.From,
			//		Content: fmt.Sprintf("%s", result.Msg),
			//	}
			//	msg, _ := json.Marshal(replyMsg)
			//	_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			//}
		}

	}
}
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Println(c.ID, "接受消息:", string(message))
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
