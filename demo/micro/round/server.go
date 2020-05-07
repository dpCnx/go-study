package main

import (
	"context"
	hello "github.com/dpCnx/go-study/demo/micro/round/proto"
	"github.com/micro/go-micro/v2"
	"log"
	"time"
)

type Say1 struct{}

func (s *Say1) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.greeter"),
		micro.RegisterTTL(time.Second*30),      //服务过期时间
		micro.RegisterInterval(time.Second*10), //注册服务的heartbeat
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	hello.RegisterSayHandler(service.Server(), new(Say1))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
