package main

import (
	"context"
	twoserver_proto "github.com/dpCnx/go-study/demo/micro/twoserver/proto"
	"github.com/micro/go-micro/v2"
	"log"
)

func main() {
	service := micro.NewService()
	service.Init()
	sTwo := twoserver_proto.NewHelloServerService("service.two", service.Client())
	if respose, err := sTwo.SayHello(context.Background(), &twoserver_proto.HelloRequest{
		Name: "DDD",
	}); err != nil {
		log.Println(err)
	} else {
		log.Println(respose.Message)
	}
}
