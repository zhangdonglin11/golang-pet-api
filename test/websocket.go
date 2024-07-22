package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang-pet-api/models/model"
	"gopkg.in/fatih/set.v0"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Chat(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer func(ws *websocket.Conn) {
	// 	err := ws.Close()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }(ws)

	iuserId, _ := strconv.Atoi(c.Query("userId"))
	userId := uint(iuserId)
	// targetId, _ := strconv.Atoi(c.Query("targetId"))

	node := Node{
		Conn:      ws,
		DataQueue: make(chan []byte),
		GroupSet:  set.New(set.ThreadSafe),
	}
	rwLocker.Lock()
	clientMap[userId] = &node
	rwLocker.Unlock()

	go sendProc(&node)
	go recvProc(&node)

	sendP2PMsg(userId, userId, []byte("hello"))

	// for {
	// 	utils.Publish(c, utils.PublishKey, "hello")
	// 	msg, err := utils.Subscribe(c, utils.PublishKey)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// }
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSet  set.Interface
}

var clientMap map[uint]*Node = make(map[uint]*Node)

var rwLocker sync.RWMutex

func sendProc(node *Node) {
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(node.Conn)

	for {
		msg := <-node.DataQueue
		fmt.Println("sendProc <-", string(msg))
		err := node.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func recvProc(node *Node) {
	for {
		_, p, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadcast(p)
	}
}

var udpsendChan = make(chan []byte)

func broadcast(p []byte) {
	udpsendChan <- p
}

func init() {
	fmt.Println("init 初始化")
	go udpSendProc()
	go udpRecvProc()
}

func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8081,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for {
		select {
		case data := <-udpsendChan:
			conn.Write(data)
		}
	}
}

func udpRecvProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 8081,
	})
	if err != nil {
		panic("udp listen err: " + err.Error())
	}
	defer conn.Close()
	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			break
		}
		dispatch(data[:n])
	}
}

func dispatch(p []byte) {
	msg := model.Message{}
	msg.CreateTime = time.Now().Unix()
	err := json.Unmarshal(p, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("接收到的数据", msg)
	switch msg.Type {
	case 1:
		sendP2PMsg(msg.UserId, msg.TargetId, p)
		//case models.GROUP:
		//	sendGroupMsg(msg.UserId, msg.TargetId, p)
		// case models.BROADCAST: sendBroadCastMsg(msg)
	}
}

func sendP2PMsg(fromId uint, targetId uint, data []byte) {
	rwLocker.RLock()
	node, ok := clientMap[targetId]
	rwLocker.RUnlock()
	// 写Redis中的消息
	//var key string
	//if fromId > targetId {
	//	key = fmt.Sprintf("msg_%d_%d", targetId, fromId)
	//} else {
	//	key = fmt.Sprintf("msg_%d_%d", fromId, targetId)
	//}
	//if _, err := utils.Rdb.ZAdd(context.Background(), key, &redis.Z{Score: float64(time.Now().Unix()), Member: data}).Result(); err != nil {
	//	fmt.Println(err)
	//}

	if !ok {
		fmt.Println("对方不在线, targetId = ", targetId)
		return
	}
	fmt.Println("sendP2PMsg", fromId, targetId, string(data))
	node.DataQueue <- data
}

//func sendGroupMsg(fromId uint, targetId uint, data []byte) {
//	// TODO: 性能不好，每次都要查群里的所有人
//	fmt.Println("sendGroupMsg", fromId, targetId, string(data))
//	contacts := models.ContactsOfCommunity(targetId)
//	for _, c := range contacts {
//		if c.OwnerId != fromId {
//			sendP2PMsg(fromId, c.OwnerId, data)
//		}
//	}
//}

func main() {
	router := gin.Default()

	router.GET("/ping", Chat)
	router.Run("127.0.0.1:8080")
}
