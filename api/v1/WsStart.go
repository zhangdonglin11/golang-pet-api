package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"golang-pet-api/global"
	"golang-pet-api/models/forms"
	"golang-pet-api/models/model"
	"golang-pet-api/utils/e"
	"log"
)

// Client 用户类
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

// 广播类，包括广播内容和源用户
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// 用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

func (manager *ClientManager) Start() {
	for {
		log.Println("<---监听管道通信--->")
		select {
		// 建立连接
		case conn := <-Manager.Register:
			log.Printf("建立新连接: %v", conn.ID)
			Manager.Clients[conn.ID] = conn
			replyMsg := &forms.ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "已连接至服务器",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
		// 断开连接
		case conn := <-Manager.Unregister:
			log.Printf("连接失败:%v", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := &forms.ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "连接已断开",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		//广播信息
		case broadcast := <-Manager.Broadcast:
			message := broadcast.Message
			sendId := broadcast.Client.SendID

			flag := false // 默认对方不在线
			//if conn, ok := Manager.Clients[sendId]; ok {
			//	select {
			//	case conn.Send <- message:
			//		flag = true
			//	default:
			//		close(conn.Send)
			//		delete(Manager.Clients, conn.ID)
			//	}
			//}
			for id, conn := range Manager.Clients {
				if id != sendId {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}

			id := broadcast.Client.ID
			replyMsg := &forms.ReplyMsg{}
			_ = json.Unmarshal(message, &replyMsg)

			if flag {
				log.Println("对方在线应答")
				replyMsg.Code = e.WebsocketOnlineReply
				replyMsg.Msg = "对方在线应答"
				msg, err := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				//将消息保存数据库中
				chatMessage := model.ChatMessage{
					ChatId:  id,
					Media:   replyMsg.Media,
					Content: replyMsg.Content,
				}
				global.Db.Create(&chatMessage)
				if err != nil {
					fmt.Println("InsertOneMsg Err", err)
				}
			} else {
				log.Println("对方不在线")
				replyMsg.Code = e.WebsocketOfflineReply
				replyMsg.Msg = "对方不在线应答"
				msg, err := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				//将消息保存数据库中
				if err != nil {
					fmt.Println("InsertOneMsg Err", err)
				}
			}
		}
	}
}
