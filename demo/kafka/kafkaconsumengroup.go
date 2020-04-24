package main

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
)

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
		log.Printf("Message topic:%q partition:%d offset:%d index:%d \n", msg.Topic, msg.Partition, msg.Offset, k.index)
		session.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	for i := 0; i < 8; i++ {
		go initComsumeres(i)
	}

	select {}
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
