package rabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"golang-pet-api/common/global"
	"golang-pet-api/common/utils/imageUtils"
	"log"
)

func ReceiveImg() {
	ch, err := global.Conn.Channel()
	failOnError(err, "Failed to open a channel")

	// 申请交换机
	err = ch.ExchangeDeclare(
		global.Config.RabbitMQ.Exchange, // name
		global.Config.RabbitMQ.Exchange, // type
		true,                            // durable
		false,                           // auto-deleted
		false,                           // internal
		false,                           // no-wait
		amqp.Table{ // arguments
			"x-delayed-type": "direct",
		})
	if err != nil {
		fmt.Println("111")
		failOnError(err, "交换机申请失败！")
		return
	}

	// 声明一个常规的队列, 其实这个也没必要声明,因为 exchange 会默认绑定一个队列
	q, err := ch.QueueDeclare(
		global.Config.RabbitMQ.Queue, // name
		true,                         // durable
		true,                         // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,                            // queue name
		global.Config.RabbitMQ.RoutingKey, // routing key
		global.Config.RabbitMQ.Exchange,   // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	// 这里监听的是 test_logs
	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	// 接收到的延时图片信息，查询redis确定图片信息是否存在redis有就删除图片
	go func() {
		for d := range msgs {
			TempImgAddr := imageUtils.TempRedisStore{}.Get(string(d.Body), true)
			imageUtils.DeleteImage(TempImgAddr)
			log.Printf("删除temp图片%s", TempImgAddr)
		}
	}()
	fmt.Println("Rabbit 图片消息队列开始监听")
	<-forever
}
