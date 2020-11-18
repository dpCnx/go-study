package main

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
)

func main() {

	//client()

	clientGroup()

	select {}
}

func client() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	//config.Consumer.Offsets.AutoCommit.Enable = true //自动提交？ 没找到同步提交offset 的方法
	// consumer
	consumer, err := sarama.NewConsumer([]string{"10.25.27.192:9092"}, config)
	if err != nil {
		log.Printf("create consumer error %v\n", err)
		return
	}
	defer consumer.Close()
	partitionList, err := consumer.Partitions("test")
	// 根据topic取到所有的分区
	if err != nil {
		log.Printf("consumer partitions error %v\n", err)
		return
	}
	for _, v := range partitionList {
		partitionConsumer, err := consumer.ConsumePartition("test", v, sarama.OffsetNewest)
		if err != nil {
			log.Printf("try create partitionConsumer error %v\n", err)
			return
		}

		// 异步从每个分区消费信息
		go func(pc sarama.PartitionConsumer) {
			for {
				select {
				case msg := <-pc.Messages():
					log.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s ,key:%s\n",
						msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value), string(msg.Key))

				case err := <-pc.Errors():
					log.Printf("err :%v\n", err)
				}
			}
		}(partitionConsumer)
	}
}

/*
	没有测试出负载均衡的作用
*/
func clientGroup() {
	for i := 0; i < 8; i++ {
		go initComsumeres(i)
	}
}

type KafkaConsumerGroupHandler struct {
	index int
}

//安装程序在新会话开始时运行，在消费品目标之前
func (k KafkaConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

//清理是在会话结束时运行的，一旦所有的消耗性laim goroutines都退出
func (k KafkaConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (k KafkaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		log.Printf("Message topic:%q partition:%d offset:%d index:%d vaule:%s\n", msg.Topic, msg.Partition, msg.Offset, k.index, string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}

func initComsumeres(i int) {
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	client, err := sarama.NewClient([]string{"10.25.27.192:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Close() }()

	group, err := sarama.NewConsumerGroupFromClient("A", client)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()

	go func() {
		for err := range group.Errors() {
			log.Println("ERROR", err)
		}
	}()

	for {
		topics := []string{"test"}
		handler := KafkaConsumerGroupHandler{
			index: i,
		}

		err := group.Consume(context.Background(), topics, handler)
		if err != nil {
			panic(err)
		}
	}
}
