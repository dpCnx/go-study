package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main() {

	//创建连接 Connection
	conn, err := amqp.Dial("amqp://dp:dp@192.168.172.128:5673/d")
	defer conn.Close()
	if err != nil {
		log.Println("amqp conn err:", err)
		return
	}
	//创建Channel
	c, err := conn.Channel()
	defer c.Close()
	if err != nil {
		log.Println("conn channel err:", err)
		return
	}

	log.Println("初始化channel successful")

	//尝试创建交换机
	/*
		参数：
		1. exchange:交换机名称
		2. type:交换机类型
		DIRECT("direct"),：定向
		FANOUT("fanout"),：扇形（广播），发送消息到每一个与之绑定队列。
		TOPIC("topic"),通配符的方式
		HEADERS("headers");参数匹配

		3. durable:是否持久化
		4. autoDelete:自动删除
		5. internal：内部使用。 一般false 	true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		6. arguments：参数
	*/

	if err = c.ExchangeDeclare(
		"demo_exchange",
		"fanout",
		true,
		false,
		false,
		//是否阻塞处理
		false,
		nil,
	); err != nil {
		log.Printf("初始化交换机失败%v", err)
		return
	}

	//创建队列 --------------------------------------->
	q, err := c.QueueDeclare(
		"demo_pubsub",
		true,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Printf("创建队列失败%v \n", err)
		return
	}

	//绑定队列到 exchange 中

	/*
		参数：
		1. queue：队列名称
		2. exchange：交换机名称
		3. routingKey：路由键，绑定规则
		如果交换机的类型为fanout ，routingKey设置为""
	*/

	err = c.QueueBind(
		q.Name,
		"",
		"demo_exchange",
		false,
		nil)


	//创建队列 --------------------------------------->
	q2, err := c.QueueDeclare(
		"demo_pubsub_2",
		true,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Printf("创建队列失败%v \n", err)
		return
	}

	//绑定队列到 exchange 中

	/*
		参数：
		1. queue：队列名称
		2. exchange：交换机名称
		3. routingKey：路由键，绑定规则
		如果交换机的类型为fanout ，routingKey设置为""
	*/

	err = c.QueueBind(
		q2.Name,
		"",
		"demo_exchange",
		false,
		nil)

	for i := 0; i <= 10; i++ {
		//调用channel 发送消息到队列中
		if err = c.Publish(
			//交换机的名字
			"demo_exchange",
			//队列名字
			"",
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
}
