package main

import "github.com/dpCnx/go-study/demo/rabbitmq/mq"

func main() {
	mq.ConsumeSimple(mq.InitChannel())
}
