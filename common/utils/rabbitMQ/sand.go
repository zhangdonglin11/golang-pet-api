package rabbitMQ

import (
	"github.com/streadway/amqp"
	"golang-pet-api/common/global"
	"log"
)

func SandTempImg(newFileName string) {
	// 将临时图片的消息发布在延时消息队列中，一定时间后定时删除
	ch, err := global.Conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 申请交换机
	err = ch.ExchangeDeclare(
		global.Config.RabbitMQ.Exchange,
		global.Config.RabbitMQ.Exchange,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-delayed-type": "direct",
		})
	if err != nil {
		failOnError(err, "交换机申请失败！")
		return
	}
	if err = ch.QueueBind(
		global.Config.RabbitMQ.Queue,
		global.Config.RabbitMQ.RoutingKey,
		global.Config.RabbitMQ.Exchange,
		false, nil); err != nil {
		failOnError(err, "绑定交换机失败！")
		return
	}
	body := newFileName
	// 将消息发送到延时队列上
	err = ch.Publish(
		global.Config.RabbitMQ.Exchange,   // exchange 这里为空则不选择 exchange
		global.Config.RabbitMQ.RoutingKey, // routing key
		false,                             // mandatory
		false,                             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Headers: map[string]interface{}{
				"x-delay": "300000", // 消息从交换机过期时间,毫秒（x-dead-message插件提供）
			},
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [temp] Sent %s", body)
}
