package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main() {

	//创建连接 Connection
	conn, err := amqp.Dial("amqp://dp:dp@192.168.172.128:5673/d")
	//defer conn.Close()
	if err != nil {
		log.Println("amqp conn err:", err)
		return
	}
	//创建Channel
	c, err := conn.Channel()
	//defer c.Close()
	if err != nil {
		log.Println("conn channel err:", err)
		return
	}

	log.Println("初始化channel successful")

	//接收消息
	msgs, err := c.Consume(
		"demo_router_2", // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		true, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Println("接收消息错误:", err)
		return
	}

	//启用协程处理消息
	go func() {
		for d := range msgs {
			log.Printf("received a message: %s ---> routingkey: %s,--->consumertag : %s \n", d.Body, d.RoutingKey, d.ConsumerTag)
		}
	}()

	select {}
}
