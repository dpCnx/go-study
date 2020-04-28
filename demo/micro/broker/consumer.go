package main

import (
	"github.com/micro/go-micro/v2/broker"
	"log"
)

var (
	topic = "go.micro.topic.d"
)

func subShare() {

	_, err := broker.Subscribe(topic, func(event broker.Event) error {

		log.Println("收到信息")
		log.Println(string(event.Message().Body))
		log.Println(event.Message().Header)

		return nil
	}, broker.Queue("D1"))

	if err != nil {
		log.Printf("sub err: %v\n", err)
	}
}

func sub() {

	_, err := broker.Subscribe(topic, func(event broker.Event) error {

		log.Println("收到信息")
		log.Println(string(event.Message().Body))
		log.Println(event.Message().Header)

		return nil
	})

	if err != nil {
		log.Printf("sub err: %v\n", err)
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

	subShare()

	select {}
}
