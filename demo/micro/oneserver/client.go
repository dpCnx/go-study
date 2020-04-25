package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"log"
)

func main() {
	s := micro.NewService()
	s.Init()
	c := s.Client()
	request := c.NewRequest("service.one", "ServiceOne.Hw", "DDD", client.WithContentType("application/json"))
	var respose string
	if err := c.Call(context.Background(), request, &respose); err != nil {
		log.Printf("c call err:%v\n", err)
		return
	} else {
		log.Println(respose)
	}
}
