package main

import (
	"github.com/Shopify/sarama"
	"log"
)

func main() {
	server()
}

//同步
func server() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	config.Producer.Return.Errors = true

	// 连接kafka同步
	client, err := sarama.NewSyncProducer([]string{"192.168.172.129:9092"}, config)
	if err != nil {
		log.Printf("newsyncproducer err:%v", err)
		return
	}
	defer client.Close()

	// 构造一个消息
	msg := &sarama.ProducerMessage{
		Topic:     "test",
		Value:     sarama.StringEncoder("hi"),
		Partition: 2,
		Key:       sarama.StringEncoder("test_key"),
	}
	partition, offset, err := client.SendMessage(msg)
	if err != nil {
		log.Printf("sendmessage err:%v", err)
		return
	}

	log.Printf("partition:%v---offset:%v", partition, offset)
}

//异步
func serverAsync() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer([]string{"192.168.172.129:9092"}, config)
	if err != nil {
		log.Printf("create producer error :%s\n", err.Error())
		return
	}

	defer producer.AsyncClose()

	// send message
	msg := &sarama.ProducerMessage{
		Topic: "test",
		Value: sarama.ByteEncoder([]byte("hello kafka")),
	}

	// send to chain
	producer.Input() <- msg

	go func() {
		for {
			select {
			case suc := <-producer.Successes():
				log.Printf("offset: %d,  timestamp: %s", suc.Offset, suc.Timestamp.String())
			case fail := <-producer.Errors():
				log.Printf("err: %s\n", fail.Err.Error())
			}
		}
	}()

	select {}
}
