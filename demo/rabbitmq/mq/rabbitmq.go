package mq

import (
	"github.com/streadway/amqp"
	"log"
)

//初始化channel
func InitChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://admin:admin@192.168.172.128:5673/demo")
	//defer conn.Close()
	if err != nil {
		log.Println("amqp conn err:", err)
		return nil
	}

	c, err := conn.Channel()
	//defer c.Close()
	if err != nil {
		log.Println("conn channel err:", err)
		return nil
	}

	log.Println("初始化channel successful")
	return c
}

//Simple 模式队列生产
func PublishSimple(c *amqp.Channel, msg string) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := c.QueueDeclare(
		"demo_test",
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
	c.Publish(
		//交换机的名字
		"",
		//队列名字
		"demo_test",
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
}

//simple 模式下消费者
func ConsumeSimple(c *amqp.Channel) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := c.QueueDeclare(
		"demo_test",
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

	//接收消息
	msgs, err := c.Consume(
		q.Name, // queue
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

	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			log.Printf("Received a message: %s", d.Body)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

//订阅模式生产
func PublishPub(c *amqp.Channel, message string) {
	//1.尝试创建交换机
	err := c.ExchangeDeclare(
		"demoexchange",
		//The common types are "direct", "fanout", "topic" ,"headers".
		"fanout",
		//是否持久化
		true,
		//是否自动删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞处理
		false,
		nil,
	)

	if err != nil {
		log.Printf("初始化交换机失败%v", err)
		return
	}

	//2.发送消息
	err = c.Publish(
		"demoexchange",
		"",
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//订阅模式消费端代码
func RecieveSub(c *amqp.Channel) {
	//1.试探性创建交换机
	err := c.ExchangeDeclare(
		"demoexchange",
		//交换机类型
		"fanout",
		true,
		false,
		//YES表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("初始化交换机失败%v", err)
		return
	}
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := c.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Printf("试探性创建队列失败%v", err)
		return
	}

	//绑定队列到 exchange 中
	err = c.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		"",
		"demoexchange",
		false,
		nil)

	//消费消息
	messges, err := c.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Println("退出请按 CTRL+C\n")
	<-forever
}

//路由模式发送消息
func PublishRouting(c *amqp.Channel, message string, key string) {
	//1.尝试创建交换机
	err := c.ExchangeDeclare(
		"demoexchange2",
		//要改成direct
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Printf("初始化交换机失败%v", err)
		return
	}

	//2.发送消息
	err = c.Publish(
		"demoexchange2",
		//要设置
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//路由模式接受消息
func RecieveRouting(c *amqp.Channel, key string) {
	//1.试探性创建交换机
	err := c.ExchangeDeclare(
		"demoexchange2",
		//交换机类型
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("初始化交换机失败%v", err)
		return
	}
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := c.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Printf("试探性创建队列失败%v", err)
		return
	}

	//绑定队列到 exchange 中
	err = c.QueueBind(
		q.Name,
		//需要绑定key
		key,
		"demoexchange2",
		false,
		nil)

	//消费消息
	messges, err := c.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Println("退出请按 CTRL+C\n")
	<-forever
}

//话题模式发送消息
func PublishTopic(c *amqp.Channel, message string, key string) {
	//1.尝试创建交换机
	err := c.ExchangeDeclare(
		"demoexchange3",
		//要改成topic
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Printf("初始化交换机失败%v", err)
		return
	}

	//2.发送消息
	err = c.Publish(
		"demoexchange3",
		//要设置
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//话题模式接受消息
//要注意key,规则
//其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
//匹配 imooc.* 表示匹配 imooc.hello, 但是imooc.hello.one需要用imooc.#才能匹配到
func RecieveTopic(c *amqp.Channel, key string) {
	//1.试探性创建交换机
	err := c.ExchangeDeclare(
		"demoexchange3",
		//交换机类型
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("初始化交换机失败%v", err)
		return
	}
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := c.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Printf("试探性创建队列失败%v", err)
		return
	}

	//绑定队列到 exchange 中
	err = c.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		key,
		"demoexchange3",
		false,
		nil)

	//消费消息
	messges, err := c.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Println("退出请按 CTRL+C\n")
	<-forever
}
