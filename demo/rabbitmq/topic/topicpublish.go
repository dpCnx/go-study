package main

import (
	"github.com/dpCnx/go-study/demo/rabbitmq/mq"
	"strconv"
)

func main() {
	c := mq.InitChannel()
	for i := 0; i < 10; i++ {
		mq.PublishTopic(c, "i am "+strconv.Itoa(i), "topic.one")
		mq.PublishTopic(c, "i am "+strconv.Itoa(i), "topic.two")
	}
}
