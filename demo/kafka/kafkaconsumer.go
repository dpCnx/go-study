package main

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"log"
)

func main() {
	var Address = []string{"10.25.27.192:9092"}
	topic := []string{"test"}

	go initComsumer("D", Address, topic)
	go initComsumer("D", Address, topic)
	go initComsumer("D", Address, topic)

	select {}

}

func initComsumer(groubid string, Address []string, topic []string) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumer, err := cluster.NewConsumer(Address, groubid, topic, config)
	if err != nil {
		log.Printf("init err :%v \n", err)
	}
	//defer consumer.Close()
	go func() {
		for {
			select {
			case msg := <-consumer.Messages():
				log.Printf("groupId=%s, topic=%s, partion=%d, offset=%d, key=%s, value=%s\n",
					groubid, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
			case err := <-consumer.Errors():
				log.Printf("err :%v\n", err)
			}

		}
	}()
}
