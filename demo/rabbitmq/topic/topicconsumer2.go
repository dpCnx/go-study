package main

import "github.com/dpCnx/go-study/demo/rabbitmq/mq"

func main() {
	c := mq.InitChannel()
	mq.RecieveTopic(c, "*.two")
}
