package main

import (
	"fmt"
	"github.com/micro/go-micro/v2/broker"
	"log"
	"time"
)

var (
	topic1 = "go.micro.topic.d"
)

func pub() {
	ticker := time.NewTicker(1 * time.Second)
	i := 0
	for _ = range ticker.C {
		msg := &broker.Message{
			Header: map[string]string{
				"id": fmt.Sprintf("%d", i),
			},
			Body: []byte(fmt.Sprintf("time %s", time.Now().String())),
		}
		if err := broker.Publish(topic1, msg); err != nil {
			log.Printf("broker publish err:%v", err)
		}

		i++
	}

}

func main() {

	if err := broker.Init(); err != nil {
		log.Printf("broker init err:%v \n", err)
		return
	}

	if err := broker.Connect(); err != nil {
		log.Printf("broker connect err:%v \n", err)
		return
	}

	pub()

	select {}
}
