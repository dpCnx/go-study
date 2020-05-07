package main

import (
	"context"
	hello "github.com/dpCnx/go-study/demo/micro/round/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-plugins/wrapper/select/roundrobin/v2"
	"log"
)

func main() {
	wrapper := roundrobin.NewClientWrapper()

	service := micro.NewService(
		micro.WrapClient(wrapper),
	)

	// parse command line flags
	service.Init()

	sTwo := hello.NewSayService("go.micro.srv.greeter", service.Client())
	if respose, err := sTwo.Hello(context.Background(), &hello.Request{
		Name: "DDD",
	}); err != nil {
		log.Println(err)
	} else {
		log.Println(respose.Msg)
	}
}
