package main

import "github.com/dpCnx/go-study/demo/rabbitmq/mq"

func main() {
	c := mq.InitChannel()
	mq.RecieveSub(c)
}
