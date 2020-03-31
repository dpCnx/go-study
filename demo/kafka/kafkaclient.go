package main

import (
	"github.com/Shopify/sarama"
	"log"
)

func main() {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	//config.Consumer.Offsets.AutoCommit.Enable = true //自动提交？ 没找到同步提交offset 的方法

	// consumer
	consumer, err := sarama.NewConsumer([]string{"192.168.60.100:9092"}, config)
	if err != nil {
		log.Printf("create consumer error %v\n", err)
		return
	}

	defer consumer.Close()

	partitionList, err := consumer.Partitions("test") // 根据topic取到所有的分区
	if err != nil {
		log.Printf("consumer partitions error %v\n", err)
		return
	}

	for _, v := range partitionList {
		partitionConsumer, err := consumer.ConsumePartition("test", v, sarama.OffsetOldest)
		if err != nil {
			log.Printf("try create partitionConsumer error %v\n", err)
			return
		}


		// 异步从每个分区消费信息
		go func(pc sarama.PartitionConsumer) {
			for {
				select {
				case msg := <-pc.Messages():
					log.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s\n",
						msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value))
				case err := <-pc.Errors():
					log.Printf("err :%v\n", err)
				}
			}
		}(partitionConsumer)
	}

}

