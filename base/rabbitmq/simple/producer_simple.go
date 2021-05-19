package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main() {

	//创建连接 Connection
	conn, err := amqp.Dial("amqp://dp:dp@192.168.172.128:5673/d")
	if err != nil {
		log.Println("amqp conn err:", err)
		return
	}
	defer conn.Close()
	//创建Channel
	c, err := conn.Channel()
	if err != nil {
		log.Println("conn channel err:", err)
		return
	}
	defer c.Close()

	log.Println("初始化channel successful")

	//创建队列Queue，如果队列不存在会自动创建，存在则跳过创建
	_, err = c.QueueDeclare(
		"demo_simple",
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		log.Println("申请队列失败:", err)
		return
	}
	//调用channel 发送消息到队列中
	if err = c.Publish(
		//交换机的名字
		"",
		//队列名字
		"demo_simple",
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("hello_demo"),
		}); err != nil {
		log.Println("申请队列失败:", err)
		return
	}
}
