package initialize

import (
	"fmt"
	"github.com/streadway/amqp"
	"golang-pet-api/global"
	"log"
)

func InitRabbitMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	global.Conn = conn
	fmt.Println("rabbitMQ 连接成功")
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
